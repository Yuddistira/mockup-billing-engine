package usecase

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (u *Usecase) SkipHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form (important for POST requests)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var (
		paramBillID int
		err         error
	)
	paramBillID, err = strconv.Atoi(r.FormValue("billing_id"))
	if err != nil {
		http.Error(w, "Bad request - parse billing_id", http.StatusBadRequest)
		return
	}

	billing, err := u.Repo.GetBilling(paramBillID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get billing data, Billing ID:%d", paramBillID), http.StatusInternalServerError)
		return
	}

	// borrower cannot skip the last payment schedule
	if billing.Tenor <= billing.CurrentPaymentIdx {
		http.Error(w, "Cannot skip last payment schedule", http.StatusInternalServerError)
		return
	}

	var amountOfSkip = billing.CurrentPaymentIdx - billing.LastPaymentIdx

	//Update Delinquent Status
	tx, err := u.Repo.BeginTx()
	if err != nil {
		http.Error(w, "Failed start transaction database", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// update payment idx
	billing.CurrentPaymentIdx = billing.CurrentPaymentIdx + 1
	err = u.Repo.UpdateBillingPaymentIdx(tx, billing.Id, billing.LastPaymentIdx, billing.CurrentPaymentIdx)
	if err != nil {
		http.Error(w, "Failed update billing payment", http.StatusInternalServerError)
		return
	}

	// update deliquent status if value deliquent is false and payment skip 2 or more than 2 times.
	if (amountOfSkip >= 2) && !billing.GetIsDelinquent() {
		err = u.Repo.UpdateBillingIsDelinquent(tx, billing.Id, true)
		if err != nil {
			http.Error(w, "Failed update deliquent data", http.StatusInternalServerError)
			return
		}
		billing.IsDelinquent = true
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		if err != nil {
			http.Error(w, "Failed commit transaction database", http.StatusInternalServerError)
			return
		}
	}

	// update template page
	billingHistories, err := u.Repo.GetAllHistBilling(billing.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get billing data, Billing ID:%d", billing.Id), http.StatusInternalServerError)
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

	if billing.GetIsDelinquent() {
		templateData.Status = "Deliquent"
	}

	templateData.Interest = templateData.Interest * amountOfSkip

	// pay the rest if tenor already reach last scheduled
	if billing.Tenor == billing.CurrentPaymentIdx {
		templateData.Interest = billing.OutstandingAmount
	} else {
		templateData.Interest = templateData.Interest * 2
	}

	t := template.Must(template.New("rows").Parse(SimulationResultPage))
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, templateData)
}
