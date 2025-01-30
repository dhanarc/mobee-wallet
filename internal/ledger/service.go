package ledger

import (
	"context"
	"math/big"
)

type Service interface {
	GetUserBalance(ctx context.Context, username string) (*big.Float, error)
	AddTransaction(ctx context.Context, username string, transaction *Transaction) error
	Transfer(ctx context.Context, request TransferRequest) error
	GetHistory(ctx context.Context, username string, page, size int) ([]Transaction, error)
}
