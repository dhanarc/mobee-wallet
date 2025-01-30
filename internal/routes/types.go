package routes

import (
	"github.com/dhanarc/mobee-wallet/internal/user/auth"
	"github.com/dhanarc/mobee-wallet/internal/wallet"
	"net/http"
	"strconv"
)

type UserRegisterRequest struct {
	Username string `json:"username"`
}

type UserRegisterResponse struct {
	Token    string `json:"token"`
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserBalanceResponse struct {
	Balance string `json:"balance"`
}

type UserTopUpRequest struct {
	Balance string `json:"balance"`
}

type WithdrawRequest struct {
	Amount string `json:"amount"`
}

type TransferRequest struct {
	ToUsername string `json:"to_username"`
	Amount     string `json:"amount"`
}

type PaginationParams struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p *PaginationParams) Parse(r *http.Request) {
	page := 1
	size := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	if s := r.URL.Query().Get("size"); s != "" {
		if val, err := strconv.Atoi(s); err == nil && val > 0 {
			size = val
		}
	}

	p.Page = page
	p.Size = size
}

type Services struct {
	AuthClient    *auth.Client
	WalletService wallet.Service
}
