package shared

import "errors"

var (
	ErrSystemError         = errors.New("system error")
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrInvalidUser         = errors.New("invalid user")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrUnauthorized        = errors.New("unauthorized")
)
