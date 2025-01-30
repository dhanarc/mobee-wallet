package routes

import (
	"context"
	"github.com/dhanarc/mobee-wallet/internal/shared"
	"github.com/dhanarc/mobee-wallet/internal/user/auth"
	"net/http"
	"strings"
)

type WalletAuthentication interface {
	Middleware(next http.Handler) http.Handler
}

type walletAuth struct {
	authClient *auth.Client
}

func NewWalletAuthentication(authClient *auth.Client) WalletAuthentication {
	return &walletAuth{
		authClient: authClient,
	}
}

func (wa *walletAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		username, err := wa.authClient.Verify(token)
		if err != nil {
			HandleError(w, r, shared.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", *username)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
