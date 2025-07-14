package repo

const (
	queryInsertMstBilling = `INSERT INTO master_billing (
    loan_amount, tenor, tenor_period,
    interest_percentage, interest_amount,
    is_delinquent, outstanding_amount,
    last_payment_idx, current_payment_idx, create_time
  ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	queryInsertHstBilling = `
  INSERT INTO history_billing (
    billing_id, payment_idx, create_time
)
VALUES (?, ?, ?)`
)

func (c *Client) InsertMasterBilling(param TableMasterBilling) error {

	_, err := c.db.Exec(queryInsertMstBilling, param.LoanAmount, param.Tenor, param.TenorPeriod, param.InterestPercentage, param.InterestAmount,
		param.IsDelinquent, param.OutstandingAmount, param.LastPaymentIdx, param.CurrentPaymentIdx, param.CreateTime)
	return err

}

func (c *Client) InsertHistoryBilling(param TableHistoryBilling) error {

	_, err := c.db.Exec(queryInsertHstBilling, param.BillingID, param.PaymentIdx, param.CreateTime)
	return err

}
