package usecase

type ScheduleBilling struct {
	PaymentIdx string
	Amount     string
	PayStatus  string
}

type SimulateBillingResp struct {
	Billings          []ScheduleBilling
	BillSchedule      string
	BillingID         int
	Status            string
	OutstandingAmount int
	Interest          int
	Finish            string
}

type PaymentRequest struct {
	BillingID int
	Amount    int
}
