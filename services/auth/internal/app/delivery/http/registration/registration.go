package registration

import (
	"auth/internal/app/models"
	"context"
	resp "github.com/R1ckNash/Bank/pkg/api/response"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
)

type UserCreator interface {
	RegisterUser(ctx context.Context, user *models.User) error
}

type Request struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
}

func New(logger *zap.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.registration.New"

		var req Request

		log := logger.With(zap.String("op", op), zap.String("username", req.Username))
		log.Info("registration new user")

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", zap.Error(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		user := &models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		err = userCreator.RegisterUser(context.Background(), user)
		if err != nil {
			log.Error("failed to register user", zap.Error(err))
			render.JSON(w, r, resp.Error("failed to register user"))
			return
		}

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, resp.OK())
}
