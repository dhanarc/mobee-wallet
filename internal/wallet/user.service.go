package wallet

import (
	"context"
	"github.com/dhanarc/mobee-wallet/internal/ledger"
	"github.com/dhanarc/mobee-wallet/internal/shared"
	"math/big"
)

func (c *client) Register(ctx context.Context, username string) (*User, *string, error) {
	createdUser, token, err := c.userService.Register(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	usr := new(User)
	usr.ID = createdUser.ID
	usr.Username = createdUser.Username
	usr.LastTransactionTime = createdUser.LastTransactionTime

	return usr, token, nil
}

func (c *client) GetBalance(ctx context.Context) (string, error) {
	currentSession, err := c.ExtractSession(ctx)
	if err != nil {
		return "", err
	}

	amount, err := c.ledgerService.GetUserBalance(ctx, *currentSession)
	if err != nil {
		return "", err
	}
	return amount.Text('f', shared.PrecisionAmount), nil
}

func (c *client) GetHistory(ctx context.Context, page, size int) ([]Transaction, error) {
	currentSession, err := c.ExtractSession(ctx)
	if err != nil {
		return nil, err
	}

	transactions, err := c.ledgerService.GetHistory(ctx, *currentSession, page, size)
	if err != nil {
		return nil, err
	}

	txns := make([]Transaction, len(transactions))
	for i, transaction := range transactions {
		txn := Transaction{
			Amount:      transaction.Amount.Text('f', shared.PrecisionAmount),
			Description: transaction.Description,
			Timestamp:   transaction.CreatedAt,
		}

		if transaction.Type == ledger.DEBIT {
			txn.Amount = new(big.Float).Mul(transaction.Amount, new(big.Float).SetFloat64(-1)).String()
		}

		txns[i] = txn
	}
	return txns, nil
}
