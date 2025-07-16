package repo

import "database/sql"

type TableMasterBilling struct {
	Id                 int          `db:"id"`
	LoanAmount         int          `db:"loan_amount"`
	Tenor              int          `db:"tenor"`
	TenorPeriod        string       `db:"tenor_period"`
	InterestPercentage int          `db:"interest_percentage"`
	InterestAmount     int          `db:"interest_amount"`
	IsDelinquent       bool         `db:"is_delinquent"`
	OutstandingAmount  int          `db:"outstanding_amount"`
	LastPaymentIdx     int          `db:"last_payment_idx"`
	CurrentPaymentIdx  int          `db:"current_payment_idx"`
	CreateTime         sql.NullTime `db:"create_time"`
	UpdateTime         sql.NullTime `db:"update_time"`
}

func (tb *TableMasterBilling) GetOutstanding() int {
	return tb.OutstandingAmount
}

func (tb *TableMasterBilling) GetIsDelinquent() bool {
	return tb.IsDelinquent
}

type TableHistoryBilling struct {
	Id         int          `db:"id"`
	BillingID  int          `db:"billing_id"`
	PaymentIdx int          `db:"payment_idx"`
	Amount     int          `db:"amount"`
	CreateTime sql.NullTime `db:"create_time"`
}
