package post_handler

import (
	"account/internal/config"
	"account/internal/models"
	slog_helper "account/internal/slog"
	"context"
	"fmt"
	resp "github.com/R1ckNash/Bank/pkg/api/response"
	"github.com/R1ckNash/Bank/pkg/middleware/auth"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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

func New(log *slog.Logger, accountCreator AccountCreator, cfg *config.Config) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		const op = "handlers.account.post.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Received request for account registration")

		// get user id from context
		userID, ok := auth.GetUserID(r)
		if !ok {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// verify user id
		verifyURL := fmt.Sprintf("http://%s:%d/user/verify/%s", cfg.AuthService.Host, cfg.AuthService.Port, userID)

		authResp, err := http.Get(verifyURL)
		if err != nil || authResp.StatusCode != 200 {
			http.Error(writer, "User not found in auth service", http.StatusForbidden)
			return
		}

		var req Request

		id, err := uuid.Parse(userID)
		if err != nil {
			log.Error("Failed to parse user id")
			http.Error(writer, `{"error": "failed to parse user id"}`, http.StatusBadRequest)
		}

		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode body")
			http.Error(writer, `{"error": "failed to decode body"}`, http.StatusBadRequest)
		}

		account := &models.Account{
			OwnerID:  id,
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
