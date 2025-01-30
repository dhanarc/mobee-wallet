package postgresql

import (
	"context"
	"errors"
	"github.com/dhanarc/mobee-wallet/internal/user"
	"gorm.io/gorm"
)

func (c *client) Register(ctx context.Context, username string) (*user.User, *string, error) {
	var newUser user.User
	err := c.db.WithContext(ctx).Where("username = ?", username).First(&newUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	if err != nil {
		newUser = user.User{
			Username: username,
		}
		if err := c.db.WithContext(ctx).Create(&newUser).Error; err != nil {
			return nil, nil, err
		}
	}

	token, err := c.auth.GenerateToken(username)
	if err != nil {
		return nil, nil, err
	}
	return &newUser, token, nil
}
