package postgresql

import (
	"github.com/dhanarc/mobee-wallet/internal/user"
	"github.com/dhanarc/mobee-wallet/internal/user/auth"
	"gorm.io/gorm"
)

type client struct {
	db   *gorm.DB
	auth *auth.Client
}

func NewClient(db *gorm.DB, auth *auth.Client) user.Service {
	return &client{
		db:   db,
		auth: auth,
	}
}
