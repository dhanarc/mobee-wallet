package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/dhanarc/mobee-wallet/internal/ledger"
	"github.com/dhanarc/mobee-wallet/internal/shared"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math/big"
	"time"
)

func (c *client) GetUserBalance(ctx context.Context, username string) (*big.Float, error) {
	var (
		err           error
		user          ledger.User
		latestAccount ledger.Account
	)

	// validate user id
	err = c.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrInvalidUser
		}
		return nil, shared.ErrSystemError
	}

	// get latest coa
	err = c.db.
		WithContext(ctx).
		Limit(1).
		Order("created_at desc").
		Take(&latestAccount, "user_id = ?", user.ID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrSystemError
	}

	// return 0 if no transaction yet
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return new(big.Float).SetFloat64(0), nil
	}

	if balance, ok := new(big.Float).SetString(latestAccount.Balance); ok {
		return balance, nil
	}
	return nil, shared.ErrSystemError
}

func (c *client) createTransaction(ctx context.Context, tx *gorm.DB, username string, transaction *ledger.Transaction) error {
	if tx == nil {
		tx = c.db.WithContext(ctx)
	}

	var txErr error
	var user ledger.User
	var latestAccount ledger.Account

	// lock user account
	txErr = tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("username = ?", username).
		First(&user).Error
	if txErr != nil {
		return txErr
	}

	// get latest balance from last account
	txErr = tx.
		Limit(1).
		Order("created_at desc").
		Take(&latestAccount, "user_id = ?", user.ID).Error
	if txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		return txErr
	}

	// calculate balance
	if errors.Is(txErr, gorm.ErrRecordNotFound) {
		transaction.CalculateBalanceAfter(new(big.Float).SetFloat64(0))
	} else {
		lastBalance, ok := new(big.Float).SetString(latestAccount.Balance)
		if !ok {
			return shared.ErrSystemError
		}
		transaction.CalculateBalanceAfter(lastBalance)
	}

	// append new coa to ledgers
	newAccount := transaction.ParseAccount(user.ID)
	txErr = tx.Create(newAccount).Error
	if txErr != nil {
		return txErr
	}

	// release lock from user record
	txErr = tx.Model(&user).Updates(ledger.User{
		LastTransactionTime: lo.ToPtr(newAccount.CreatedAt),
	}).Error
	if txErr != nil {
		return txErr
	}

	return nil
}

func (c *client) AddTransaction(ctx context.Context, username string, transaction *ledger.Transaction) error {
	if err := transaction.Validate(); err != nil {
		return err
	}

	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return c.createTransaction(ctx, tx, username, transaction)
	})
}

func (c *client) doTransfer(ctx context.Context, tx *gorm.DB, request ledger.TransferRequest) error {
	var err error
	if tx == nil {
		tx = c.db.WithContext(ctx).Begin()
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	defer func(err error) {
		if err != nil {
			tx.Rollback()
		}
	}(err)

	now := time.Now()

	err = c.createTransaction(ctx, tx, request.Origin, &ledger.Transaction{
		Date:        now,
		Type:        ledger.DEBIT,
		Amount:      request.Amount,
		Description: fmt.Sprintf("Transfer amount %s to %s", request.Amount.Text('f', shared.PrecisionAmount), request.Destination),
	})
	if err != nil {
		return err
	}

	err = c.createTransaction(ctx, tx, request.Destination, &ledger.Transaction{
		Date:        now,
		Type:        ledger.CREDIT,
		Amount:      request.Amount,
		Description: fmt.Sprintf("Transfer amount %s from %s", request.Amount.Text('f', shared.PrecisionAmount), request.Origin),
	})
	if err != nil {
		return err
	}

	return tx.Commit().Error
}

func (c *client) Transfer(ctx context.Context, request ledger.TransferRequest) error {
	var err error

	// get current user balance
	userBalance, err := c.GetUserBalance(ctx, request.Origin)
	if err != nil {
		return err
	}

	// validate balance and amount transfer
	if userBalance == nil {
		return shared.ErrSystemError
	}

	if userBalance.Cmp(request.Amount) < 0 {
		return shared.ErrInsufficientBalance
	}

	return c.doTransfer(ctx, nil, request)
}

func (c *client) GetHistory(ctx context.Context, username string, page, size int) ([]ledger.Transaction, error) {
	var err error
	var user ledger.User
	var accounts []ledger.Account

	// validate user id
	err = c.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrInvalidUser
		}
		return nil, shared.ErrSystemError
	}

	// get latest coa
	offset := (page - 1) * size
	err = c.db.
		WithContext(ctx).
		Limit(size).
		Offset(offset).
		Order("created_at desc").
		Take(&accounts, "user_id = ?", user.ID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.ErrSystemError
	}

	transactions := make([]ledger.Transaction, len(accounts))
	for i, account := range accounts {
		amount, ok := new(big.Float).SetString(account.Amount)
		if !ok {
			return nil, shared.ErrSystemError
		}

		balance, ok := new(big.Float).SetString(account.Balance)
		if !ok {
			return nil, shared.ErrSystemError
		}
		transactions[i] = ledger.Transaction{
			Date:         account.TransactionDate,
			Type:         account.TransactionType,
			Amount:       amount,
			BalanceAfter: balance,
			Description:  account.Description,
			CreatedAt:    account.CreatedAt,
		}
	}
	return transactions, nil
}
