package wallet

import (
	"context"
	"fmt"
	"github.com/dhanarc/mobee-wallet/internal/ledger"
	"github.com/dhanarc/mobee-wallet/internal/shared"
	"math/big"
	"time"
)

func (c *client) Deposit(ctx context.Context, amount string) error {
	currentSession, err := c.ExtractSession(ctx)
	if err != nil {
		return err
	}

	txAmount, ok := new(big.Float).SetString(amount)
	if !ok {
		return shared.ErrInvalidArgument
	}

	if txAmount.Cmp(new(big.Float).SetFloat64(0)) < 0 {
		return shared.ErrInvalidArgument
	}

	now := time.Now()
	return c.ledgerService.AddTransaction(ctx, *currentSession, &ledger.Transaction{
		Date:        now,
		Type:        ledger.CREDIT,
		Amount:      txAmount,
		Description: fmt.Sprintf("Topup Balance %s", txAmount.Text('f', shared.PrecisionAmount)),
	})
}

func (c *client) Withdraw(ctx context.Context, amount string) error {
	currentSession, err := c.ExtractSession(ctx)
	if err != nil {
		return err
	}

	txAmount, ok := new(big.Float).SetString(amount)
	if !ok {
		return shared.ErrInvalidArgument
	}

	now := time.Now()
	return c.ledgerService.AddTransaction(ctx, *currentSession, &ledger.Transaction{
		Date:        now,
		Type:        ledger.DEBIT,
		Amount:      txAmount,
		Description: fmt.Sprintf("Withdraw Balance %s", txAmount.Text('f', shared.PrecisionAmount)),
	})
}

func (c *client) Transfer(ctx context.Context, destination string, amount string) error {
	currentSession, err := c.ExtractSession(ctx)
	if err != nil {
		return err
	}

	txAmount, ok := new(big.Float).SetString(amount)
	if !ok {
		return shared.ErrInvalidArgument
	}

	return c.ledgerService.Transfer(ctx, ledger.TransferRequest{
		Origin:      *currentSession,
		Destination: destination,
		Amount:      txAmount,
	})
}
