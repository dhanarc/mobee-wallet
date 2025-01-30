package routes

import (
	"errors"
	"github.com/dhanarc/mobee-wallet/internal/shared"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func Init(services *Services) *chi.Mux {
	r := chi.NewRouter()

	handlers := NewHandler(services.WalletService)

	r.Use(middleware.Logger)
	r.Post("/register", handlers.RegisterUser)

	r.Route("/", func(pr chi.Router) {
		authMiddleware := NewWalletAuthentication(services.AuthClient)
		pr.Use(authMiddleware.Middleware)

		pr.Get("/user/balance", handlers.GetBalance)
		pr.Get("/user/history", handlers.GetHistory)

		pr.Post("/transaction/deposit", handlers.Deposit)
		pr.Post("/transaction/withdraw", handlers.Withdraw)
		pr.Post("/transaction/transfer", handlers.Transfer)
	})

	return r
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	httpCode := http.StatusInternalServerError
	msg := err.Error()
	defer func() {
		render.Status(r, httpCode)
		render.PlainText(w, r, msg)
	}()

	switch {
	case errors.Is(err, shared.ErrInsufficientBalance) || errors.Is(err, shared.ErrInvalidArgument):
		httpCode = http.StatusBadRequest
	case errors.Is(err, shared.ErrUnauthorized):
		httpCode = http.StatusUnauthorized
	case errors.Is(err, shared.ErrInvalidUser):
		httpCode = http.StatusNotFound
	default:
		httpCode = http.StatusInternalServerError
	}
}
