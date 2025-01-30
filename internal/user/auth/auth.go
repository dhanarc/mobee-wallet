package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"time"
)

type Client struct {
	key string
}

func NewClient(key string) *Client {
	return &Client{
		key: key,
	}
}

func (c *Client) Verify(tokenString string) (*string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Validate the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		subject, err := claims.GetSubject()
		if err != nil {
			return nil, err
		}
		return lo.ToPtr(subject), nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (c *Client) GenerateToken(username string) (*string, error) {
	key := []byte(c.key)
	currentTime := time.Now()
	expiredTime := currentTime.Add(33 * time.Hour)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    "mobee-wallet",
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expiredTime),
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ID:        uuid.NewString(),
	})

	token, err := jwtToken.SignedString(key)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
