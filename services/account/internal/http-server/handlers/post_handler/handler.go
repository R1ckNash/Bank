package post_handler

import (
	"account/internal/models"
	slog_helper "account/internal/slog"
	"context"
	resp "github.com/R1ckNash/Bank/pkg/api/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type AccountCreator interface {
	RegisterAccount(ctx context.Context, account *models.Account) error
}

type Request struct {
	OwnerID  string `json:"owner_id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Currency string `json:"currency"`
	Email    string `json:"email"`
}

func New(log *slog.Logger, accountCreator AccountCreator) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		const op = "handlers.account.post.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Received request for account registration")

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode body")
			http.Error(writer, `{"error": "failed to decode body"}`, http.StatusBadRequest)
		}

		account := &models.Account{
			OwnerID:  req.OwnerID,
			Name:     req.Name,
			Currency: req.Currency,
			Email:    req.Email,
		}

		err = accountCreator.RegisterAccount(r.Context(), account)
		if err != nil {
			log.Error("failed to create account", slog_helper.Err(err))
			http.Error(writer, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}

		responseOK(writer, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, resp.OK())
}
