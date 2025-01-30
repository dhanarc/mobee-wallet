package postgresql

import (
	"github.com/dhanarc/mobee-wallet/internal/ledger"
	"gorm.io/gorm"
)

type client struct {
	db *gorm.DB
}

func NewClient(db *gorm.DB) ledger.Service {
	return &client{
		db: db,
	}
}
