package postgresql

import (
	"github.com/dhanarc/mobee-wallet/internal/ledger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) LockUser(tx *gorm.DB, username string) (*ledger.User, error) {
	var user ledger.User

	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
