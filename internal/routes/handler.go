package routes

import (
	"encoding/json"
	"github.com/dhanarc/mobee-wallet/internal/wallet"
	"net/http"

	"github.com/go-chi/render"
)

type Handler struct {
	walletService wallet.Service
}

func NewHandler(walletService wallet.Service) *Handler {
	return &Handler{
		walletService: walletService,
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request *UserRegisterRequest
	var err error

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	createdUser, jwtToken, err := h.walletService.Register(r.Context(), request.Username)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	response := UserRegisterResponse{
		Token:    *jwtToken,
		ID:       createdUser.ID,
		Username: createdUser.Username,
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := h.walletService.GetBalance(r.Context())
	if err != nil {
		HandleError(w, r, err)
		return
	}

	response := UserBalanceResponse{
		Balance: balance,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	qparams := new(PaginationParams)
	qparams.Parse(r)

	histories, err := h.walletService.GetHistory(r.Context(), qparams.Page, qparams.Size)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, histories)
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	var request *UserTopUpRequest
	var err error

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = h.walletService.Deposit(r.Context(), request.Balance)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.PlainText(w, r, "Deposit successful")
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	var request *WithdrawRequest
	var err error

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = h.walletService.Withdraw(r.Context(), request.Amount)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.PlainText(w, r, "Withdraw successful")
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	var request *TransferRequest
	var err error

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	err = h.walletService.Transfer(r.Context(), request.ToUsername, request.Amount)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.PlainText(w, r, "Transfer successful")
}
