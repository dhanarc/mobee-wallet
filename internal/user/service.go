package user

import "context"

type Service interface {
	Register(ctx context.Context, username string) (*User, *string, error)
}
