package login

import (
	"context"
	resp "github.com/R1ckNash/Bank/pkg/api/response"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"net/http"
)

type UserAuthenticator interface {
	LoginUser(ctx context.Context, username, password string) (string, error)
}

type Request struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" binding:"required,min=4"`
}

func New(logger *zap.Logger, userAuthenticator UserAuthenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.login.New"

		var req Request

		log := logger.With(zap.String("op", op), zap.String("username", req.Username))
		log.Info("login user")

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", zap.Error(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		token, err := userAuthenticator.LoginUser(context.Background(), req.Username, req.Password)
		if err != nil {
			log.Error("failed to decode request body", zap.Error(err))
			render.JSON(w, r, resp.Error("failed to login"))
			return
		}

		responseToken(w, r, token)
	}
}

func responseToken(w http.ResponseWriter, r *http.Request, token string) {
	render.JSON(w, r, resp.OKWithData(map[string]interface{}{
		"token": token,
	}))
}
