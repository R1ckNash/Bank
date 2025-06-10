package account

import (
	"account/internal/models"
	"account/internal/repository/account_storage"
	slog_helper "account/internal/slog"
	"context"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
	"github.com/R1ckNash/Bank/pkg/helpers"
	"github.com/R1ckNash/Bank/pkg/transaction_manager"
	"log/slog"
	"time"
)

func (s *accountService) RegisterAccount(ctx context.Context, acc *models.Account) error {
	const op = "accountService.RegisterAccount"

	log := s.Logger.With(
		slog.String("op", op),
	)

	log.Info("Processing request for create account", slog.String("owner_id", acc.OwnerID))

	//todo add validation
	//todo add idempotency

	accountDTO := &account_storage.Account{
		OwnerID:   acc.OwnerID,
		Name:      acc.Name,
		Currency:  acc.Currency,
		Email:     acc.Email,
		IsBlocked: false,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := helpers.WithRetries(ctx, func(ctx context.Context) error {
		var err error
		err = s.TransactionManager.RunReadCommitted(ctx, transaction_manager.ReadWrite,
			func(txCtx context.Context) error { // TRANSANCTION SCOPE
				if err = s.AccountStorage.CreateAccount(txCtx, accountDTO); err != nil {
					return err
				}
				return nil
			},
		)
		return err
	})

	if err != nil {
		log.Warn("error saving account", slog_helper.Err(err))
		return pkgerrors.Wrap(op, err)
	}

	return nil
}
