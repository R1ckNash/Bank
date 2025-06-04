package account

import (
	"account/internal/models"
	slog_helper "account/internal/slog"
	"context"
	"fmt"
	"log/slog"
)

func (s *accountService) GetAccount(ctx context.Context, accountID int64) (models.Account, error) {
	const op = "accountService.GetAccount"

	log := s.Logger.With(
		slog.String("op", op),
	)

	accDB, err := s.AccountStorage.GetByID(ctx, accountID)
	if err != nil {
		log.Error("failed to retrieve accDB by id", slog_helper.Err(err))
		return models.Account{}, fmt.Errorf("%s: %w", op, err)
	}

	acc := models.Account{
		OwnerID:  accDB.OwnerID,
		Name:     accDB.Name,
		Currency: accDB.Currency,
		Email:    accDB.Email,
	}

	return acc, nil
}
