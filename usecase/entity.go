package usecase

type ScheduleBilling struct {
	PaymentIdx string
	Amount     string
	PayStatus  string
}

type SimulateBillingResp struct {
	Billings          []ScheduleBilling
	Status            string
	OutstandingAmount int
	Interest          int
}
