package repo

import (
	"database/sql"
	"time"
)

const (
	queryInsertMstBilling = `INSERT INTO master_billing (
    loan_amount, tenor, tenor_period,
    interest_percentage, interest_amount,
    is_delinquent, outstanding_amount,
    last_payment_idx, current_payment_idx, create_time
  ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	queryInsertHstBilling = `
  	INSERT INTO history_billing (
    billing_id, payment_idx, amount, create_time)
	VALUES (?, ?, ?, ?)`

	queryUpdateMstBillingIsDelinquent = `
	UPDATE master_billing SET is_delinquent = ?,update_time = ? WHERE id = ?;`

	queryUpdateMstBillingOutstandingAmount = `
	UPDATE master_billing SET outstanding_amount = ?,update_time = ? WHERE id = ?;`

	queryUpdateMstBillingPaymentIdx = `
	UPDATE master_billing SET last_payment_idx = ?,current_payment_idx= ?,update_time = ? WHERE id = ?;`

	queryGetMstBilling = `Select id,loan_amount, tenor, tenor_period,interest_percentage, interest_amount,is_delinquent, outstanding_amount,
    last_payment_idx, current_payment_idx from master_billing where id = ?`

	queryGetAllHstBilling = `Select id,billing_id, payment_idx, amount,create_time from history_billing where billing_id = ?`
)

func (c *Client) InsertMasterBilling(param TableMasterBilling) (int64, error) {

	result, err := c.db.Exec(queryInsertMstBilling, param.LoanAmount, param.Tenor, param.TenorPeriod, param.InterestPercentage, param.InterestAmount,
		param.IsDelinquent, param.OutstandingAmount, param.LastPaymentIdx, param.CurrentPaymentIdx, param.CreateTime)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()

}

func (c *Client) InsertHistoryBilling(tx *sql.Tx, param TableHistoryBilling) (int64, error) {
	if tx != nil {
		_, err := tx.Exec(queryInsertHstBilling, param.BillingID, param.PaymentIdx, param.Amount, param.CreateTime)
		return 0, err
	}
	result, err := c.db.Exec(queryInsertHstBilling, param.BillingID, param.PaymentIdx, param.Amount, param.CreateTime)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()

}

func (c *Client) UpdateBillingIsDelinquent(tx *sql.Tx, billingID int, IsDelinquent bool) error {
	if tx != nil {
		_, err := tx.Exec(queryUpdateMstBillingIsDelinquent, IsDelinquent, time.Now(), billingID)
		return err
	}

	_, err := c.db.Exec(queryUpdateMstBillingIsDelinquent, IsDelinquent, time.Now(), billingID)
	return err

}

func (c *Client) UpdateBillingOutstandingAmount(tx *sql.Tx, billingID int, outstandingAmount int) error {
	if tx != nil {
		_, err := tx.Exec(queryUpdateMstBillingOutstandingAmount, outstandingAmount, time.Now(), billingID)
		return err
	}

	_, err := c.db.Exec(queryUpdateMstBillingOutstandingAmount, outstandingAmount, time.Now(), billingID)
	return err
}

func (c *Client) UpdateBillingPaymentIdx(tx *sql.Tx, billingID int, lastPaymentIdx, currentPaymentIdx int) error {
	if tx != nil {
		_, err := tx.Exec(queryUpdateMstBillingPaymentIdx, lastPaymentIdx, currentPaymentIdx, time.Now(), billingID)
		return err
	}

	_, err := c.db.Exec(queryUpdateMstBillingPaymentIdx, lastPaymentIdx, currentPaymentIdx, time.Now(), billingID)
	return err
}

func (c *Client) GetBilling(billingID int) (TableMasterBilling, error) {
	var result TableMasterBilling
	err := c.db.Get(&result, queryGetMstBilling, billingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}
		return result, err
	}

	return result, nil
}
func (c *Client) GetAllHistBilling(billingID int) ([]TableHistoryBilling, error) {

	rows, err := c.db.Query(queryGetAllHstBilling, billingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TableHistoryBilling

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var thb TableHistoryBilling
		if err := rows.Scan(&thb.Id, &thb.BillingID, &thb.PaymentIdx, &thb.Amount,
			&thb.CreateTime); err != nil {
			return result, err
		}
		result = append(result, thb)
	}
	if err = rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}
