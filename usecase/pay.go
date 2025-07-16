package usecase

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/mockup-billing-engine/repo"
)

func (u *Usecase) MakePayment(w http.ResponseWriter, r *http.Request) {
	// Parse form (important for POST requests)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var (
		request PaymentRequest
		err     error
	)
	request.BillingID, err = strconv.Atoi(r.FormValue("billing_id"))
	if err != nil {
		http.Error(w, "Bad request - parse billing_id", http.StatusBadRequest)
		return
	}
	request.Amount, err = strconv.Atoi(r.FormValue("interest"))
	if err != nil {
		http.Error(w, "Bad request - parse payment_amount", http.StatusBadRequest)
		return
	}

	if request.Amount == 0 {
		http.Error(w, "Bad request - parse payment_amount", http.StatusBadRequest)
		return
	}

	billing, err := u.Repo.GetBilling(request.BillingID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get billing data, Billing ID:%d", request.BillingID), http.StatusInternalServerError)
		return
	}

	//Update Delinquent Status
	tx, err := u.Repo.BeginTx()
	if err != nil {
		http.Error(w, "Failed start transaction database", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert multiple payment if there skip
	i := billing.LastPaymentIdx
	tempOutstanding := billing.OutstandingAmount
	for i < billing.CurrentPaymentIdx {
		// use the interest amount for payment or skip payment before, but use outstanding amount once this payment reach end of schedule.
		tempAmount := billing.InterestAmount
		if (i + 1) == billing.Tenor {
			tempAmount = tempOutstanding
		}

		_, err := u.Repo.InsertHistoryBilling(tx, repo.TableHistoryBilling{
			BillingID:  request.BillingID,
			PaymentIdx: i + 1,
			Amount:     tempAmount,
			CreateTime: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			http.Error(w, "Failed insert billing history", http.StatusInternalServerError)
			return
		}
		i++
		tempOutstanding = tempOutstanding - tempAmount
	}

	// update payment idx
	billing.LastPaymentIdx = i
	billing.CurrentPaymentIdx = billing.CurrentPaymentIdx + 1
	err = u.Repo.UpdateBillingPaymentIdx(tx, request.BillingID, billing.LastPaymentIdx, billing.CurrentPaymentIdx)
	if err != nil {
		http.Error(w, "Failed update billing payment", http.StatusInternalServerError)
		return
	}

	// update deliquent status if before is true
	if billing.GetIsDelinquent() {
		err = u.Repo.UpdateBillingIsDelinquent(tx, request.BillingID, false)
		if err != nil {
			http.Error(w, "Failed update deliquent data", http.StatusInternalServerError)
			return
		}
	}

	// update outstanding amount
	billing.OutstandingAmount = billing.OutstandingAmount - request.Amount
	err = u.Repo.UpdateBillingOutstandingAmount(tx, request.BillingID, billing.OutstandingAmount)
	if err != nil {
		http.Error(w, "Failed update outstanding amount data", http.StatusInternalServerError)
		return
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		http.Error(w, "Failed commit transaction database", http.StatusInternalServerError)
		return
	}

	// update template page
	billingHistories, err := u.Repo.GetAllHistBilling(request.BillingID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get billing data, Billing ID:%d", request.BillingID), http.StatusInternalServerError)
		return
	}

	scheduleBills := u.buildScheduleBilling(billing, billingHistories)

	templateData := SimulateBillingResp{
		Billings:          scheduleBills,
		BillSchedule:      fmt.Sprintf("%s %d", billing.TenorPeriod, billing.CurrentPaymentIdx),
		BillingID:         billing.Id,
		Status:            "Normal",
		OutstandingAmount: billing.OutstandingAmount,
		Interest:          billing.InterestAmount,
	}

	// pay the rest if tenor already reach last scheduled
	if billing.Tenor == billing.CurrentPaymentIdx {
		templateData.Interest = billing.OutstandingAmount
	}

	// no more interest to be pay if outstanding amount is 0
	if billing.OutstandingAmount == 0 {
		templateData.Interest = 0
		templateData.BillSchedule = ""
		templateData.Finish = "Finished"
	}

	t := template.Must(template.New("rows").Parse(SimulationResultPage))
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, templateData)
}
