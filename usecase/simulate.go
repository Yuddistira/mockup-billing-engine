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

func (u *Usecase) SimulateHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	// Parse form (important for POST requests)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	param := repo.TableMasterBilling{
		IsDelinquent:      false,
		LastPaymentIdx:    0,
		CurrentPaymentIdx: 1,
		CreateTime:        sql.NullTime{Time: time.Now(), Valid: true},
	}

	param.LoanAmount, err = strconv.Atoi(r.FormValue("loan"))
	if err != nil {
		http.Error(w, "Bad request - parse loan amount", http.StatusBadRequest)
		return
	}
	param.Tenor, err = strconv.Atoi(r.FormValue("tenor"))
	if err != nil {
		http.Error(w, "Bad request - parse tenor", http.StatusBadRequest)
		return
	}
	param.TenorPeriod = r.FormValue("period")
	param.InterestPercentage, err = strconv.Atoi(r.FormValue("interest"))
	if err != nil {
		http.Error(w, "Bad request - parse interest", http.StatusBadRequest)
		return
	}

	param.OutstandingAmount = int(((float32(param.LoanAmount) / 100) * float32(param.InterestPercentage)) + float32(param.LoanAmount))
	param.InterestAmount = param.OutstandingAmount / param.Tenor

	billingId, err := u.Repo.InsertMasterBilling(param)
	if err != nil {
		// http.Error(w, fmt.Sprintf("Failed Store data billing - %s", err.Error()), http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed Store data billing - %s", err.Error())
		return
	}

	scheduleBills := u.buildScheduleBilling(param, []repo.TableHistoryBilling{})

	templateData := SimulateBillingResp{
		Billings:          scheduleBills,
		BillSchedule:      fmt.Sprintf("%s %d", param.TenorPeriod, param.CurrentPaymentIdx),
		BillingID:         int(billingId),
		Status:            "Normal",
		OutstandingAmount: param.OutstandingAmount,
		Interest:          param.InterestAmount,
	}

	t := template.Must(template.New("rows").Parse(SimulationResultPage))
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, templateData)
}

func (u *Usecase) buildScheduleBilling(billingData repo.TableMasterBilling, billingHistories []repo.TableHistoryBilling) []ScheduleBilling {
	var tempScheduleBills []ScheduleBilling

	// Append paid billing to scheduled list data
	for i, v := range billingHistories {
		sb := ScheduleBilling{
			PaymentIdx: fmt.Sprintf("%s %d", billingData.TenorPeriod, v.PaymentIdx),
			Amount:     fmt.Sprintf("Rp. %d", v.Amount),
		}

		if i > billingData.LastPaymentIdx && i < billingData.CurrentPaymentIdx {
			sb.PayStatus = "Skip"
		} else {
			sb.PayStatus = "Paid"
		}

		tempScheduleBills = append(tempScheduleBills, sb)
	}

	// Append unpaid billing to scheduled list data
	idx := len(tempScheduleBills) + 1

	for idx <= billingData.Tenor {
		sb := ScheduleBilling{
			PaymentIdx: fmt.Sprintf("%s %d", billingData.TenorPeriod, idx),
			Amount:     "",
			PayStatus:  "",
		}

		if idx > billingData.LastPaymentIdx && idx < billingData.CurrentPaymentIdx {
			sb.PayStatus = "Skip"
		}

		tempScheduleBills = append(tempScheduleBills, sb)

		idx++
	}
	return tempScheduleBills
}
