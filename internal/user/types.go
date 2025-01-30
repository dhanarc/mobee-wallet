package user

import "time"

type User struct {
	ID                  int64     `gorm:"column:id" json:"id"`
	Username            string    `gorm:"column:username" json:"username"`
	LastTransactionTime time.Time `gorm:"column:last_transaction" json:"last_transaction"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"created_at"`
}
