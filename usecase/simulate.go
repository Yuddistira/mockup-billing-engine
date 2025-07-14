package usecase

import (
	"database/sql"
	"fmt"
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

	err = u.Repo.InsertMasterBilling(param)
	if err != nil {
		// http.Error(w, fmt.Sprintf("Failed Store data billing - %s", err.Error()), http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed Store data billing - %s", err.Error())
		return
	}

	w.Write([]byte(SimulationResultPage))
}
