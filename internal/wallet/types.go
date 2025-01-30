package wallet

import "time"

type Transaction struct {
	Amount      string    `json:"amount"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type User struct {
	ID                  int64     `json:"id"`
	Username            string    `json:"username"`
	LastTransactionTime time.Time `json:"last_transaction_time"`
}
