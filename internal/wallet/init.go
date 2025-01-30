package wallet

import (
	"context"
	"github.com/dhanarc/mobee-wallet/internal/ledger"
	"github.com/dhanarc/mobee-wallet/internal/shared"
	"github.com/dhanarc/mobee-wallet/internal/user"
)

type client struct {
	ledgerService ledger.Service
	userService   user.Service
}

func NewClient(ledgerService ledger.Service, userService user.Service) Service {
	return &client{
		ledgerService: ledgerService,
		userService:   userService,
	}
}

func (c *client) ExtractSession(ctx context.Context) (*string, error) {
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return nil, shared.ErrUnauthorized
	}
	return &username, nil
}
