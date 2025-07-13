package repo

import "database/sql"

const (
	createTableBilling = `CREATE TABLE IF NOT EXISTS master_billing (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    loan_amount INTEGER NOT NULL,
    tenor INTEGER NOT NULL,
    tenor_period TEXT NOT NULL,
    interest_percentage INTEGER NOT NULL,
    interest_amount INTEGER NOT NULL,
    is_delinquent INTEGER NOT NULL DEFAULT 0, -- boolean stored as 0 or 1
    outstanding_amount INTEGER,
    last_payment_idx INTEGER DEFAULT 0,
    current_payment_idx INTEGER DEFAULT 0,
    create_time DATETIME NOT NULL
);
`
	createTableHistoryBilling = `CREATE TABLE IF NOT EXISTS history_billing (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    billing_id INTEGER NOT NULL,
    payment_idx INTEGER NOT NULL,
    create_time DATETIME NOT NULL,
    FOREIGN KEY (billing_id) REFERENCES master_billing(id)
);

	`
)

func Init() {
	// Connect to or create SQLite DB file
	db, err := sql.Open("sqlite3", "./billing.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
    )`)
	if err != nil {
		panic(err)
	}

}
