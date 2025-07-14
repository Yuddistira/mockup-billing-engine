package repo

import "database/sql"

type TableMasterBilling struct {
	LoanAmount         int
	Tenor              int
	TenorPeriod        string
	InterestPercentage int
	InterestAmount     int
	IsDelinquent       bool
	OutstandingAmount  int
	LastPaymentIdx     int
	CurrentPaymentIdx  int
	CreateTime         sql.NullTime
	UpdateTime         sql.NullTime
}

type TableHistoryBilling struct {
	Id         int
	BillingID  int
	PaymentIdx int
	CreateTime sql.NullTime
}
