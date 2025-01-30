package wallet

import "context"

type Service interface {
	Register(ctx context.Context, username string) (*User, *string, error)

	Deposit(ctx context.Context, amount string) error
	Withdraw(ctx context.Context, amount string) error
	Transfer(ctx context.Context, destination string, amount string) error

	GetBalance(ctx context.Context) (string, error)
	GetHistory(ctx context.Context, page, size int) ([]Transaction, error)
}
