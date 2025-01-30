package ledger

import (
	"encoding/json"
	"github.com/dhanarc/mobee-wallet/internal/shared"
	"math/big"
	"time"
)

type TransactionType string

var (
	CREDIT TransactionType = "credit"
	DEBIT  TransactionType = "debit"
)

// data access object

type User struct {
	ID                  int64      `gorm:"column:id" json:"id"`
	Username            string     `gorm:"column:username" json:"username"`
	LastTransactionTime *time.Time `gorm:"column:last_transaction" json:"last_transaction"`
	CreatedAt           time.Time  `gorm:"column:created_at" json:"created_at"`
}

type Account struct {
	ID              int64           `gorm:"column:id" json:"id"`
	UserID          int64           `gorm:"column:user_id" json:"user_id"`
	TransactionDate time.Time       `gorm:"column:transaction_date" json:"transaction_date"`
	TransactionType TransactionType `gorm:"column:transaction_type" json:"transaction_type"`
	Amount          string          `gorm:"column:amount" json:"amount"`
	Balance         string          `gorm:"column:balance" json:"balance"`
	Description     string          `gorm:"column:description" json:"description"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"created_at"`
}

// data transfer object

type Transaction struct {
	Date         time.Time       `json:"date"`
	Type         TransactionType `json:"type"`
	Amount       *big.Float      `json:"amount"`
	BalanceAfter *big.Float      `json:"balance_after"`
	Description  string          `json:"description"`
	CreatedAt    time.Time       `json:"created_at"`
}

func (t *Transaction) ParseAccount(userID int64) *Account {
	return &Account{
		UserID:          userID,
		TransactionDate: t.Date,
		TransactionType: t.Type,
		Amount:          t.Amount.Text('f', shared.PrecisionAmount),
		Balance:         t.BalanceAfter.Text('f', shared.PrecisionAmount),
		Description:     t.Description,
		CreatedAt:       time.Now(),
	}
}

func (t *Transaction) Validate() error {
	if t.Type != CREDIT && t.Type != DEBIT {
		return shared.ErrSystemError
	}
	if t.Amount.Cmp(new(big.Float).SetFloat64(0)) < 0 {
		return shared.ErrInvalidArgument
	}

	return nil
}

func (t *Transaction) CalculateBalanceAfter(previousBalance *big.Float) {
	if t.Type == CREDIT {
		t.BalanceAfter = new(big.Float).Add(previousBalance, t.Amount)
		return
	}

	t.BalanceAfter = new(big.Float).Sub(previousBalance, t.Amount)
}

func (t *Transaction) ToJSONString() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type TransferRequest struct {
	Origin      string     `json:"origin"`
	Destination string     `json:"destination"`
	Amount      *big.Float `json:"amount"`
	Description string     `json:"description"`
}
