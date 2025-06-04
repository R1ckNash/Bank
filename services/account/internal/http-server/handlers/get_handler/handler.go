package get_handler

import (
	"account/internal/models"
	slog_helper "account/internal/slog"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type AccountGetter interface {
	GetAccount(ctx context.Context, accountID int64) (models.Account, error)
}

func New(log *slog.Logger, accountGetter AccountGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.account.get.New"

		log = log.With(
			slog.String("op", op),
		)
		var requestParam = chi.URLParam(request, "accountId")
		id, err := strconv.Atoi(requestParam)
		if err != nil {
			log.Error("failed to decode id", slog_helper.Err(err))
			http.Error(writer, `{"error": "incorrect accountId"}`, http.StatusBadRequest)
		}

		account, err := accountGetter.GetAccount(request.Context(), int64(id))
		if err != nil {
			log.Error("failed to retrieve account", slog.Int("accountId", id))
			http.Error(writer, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}

		render.JSON(writer, request, account)
	}
}
