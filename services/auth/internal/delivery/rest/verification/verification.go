package verification

import (
	"context"
	resp "github.com/R1ckNash/Bank/pkg/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

type UserVerificator interface {
	VerifyUser(ctx context.Context, id uuid.UUID) bool
}

type Request struct {
	ID string `json:"id" validate:"required"`
}

func New(logger *zap.Logger, userVerificator UserVerificator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := chi.URLParam(r, "user_id")
		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, `{"error": "invalid user_id"}`, http.StatusBadRequest)
			return
		}

		isExist := userVerificator.VerifyUser(r.Context(), userUUID)
		if !isExist {
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
			return
		}

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, resp.OK())
}
