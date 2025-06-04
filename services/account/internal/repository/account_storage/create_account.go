package account_storage

import (
	"account/internal/models"
	"context"
	"errors"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *AccountStorage) CreateAccount(ctx context.Context, acc *Account) error {
	const api = "account_storage.CreateAccount"

	query := `insert into account (owner_id, name, currency, email, is_blocked, balance, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	if _, err := s.driver.GetQueryEngine(ctx).Exec(ctx, query, acc.OwnerID, acc.Name, acc.Currency, acc.Email, acc.IsBlocked, acc.Balance, acc.CreatedAt, acc.UpdatedAt); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return pkgerrors.Wrap(api, models.ErrAlreadyExists)
		}
		return pkgerrors.Wrap(api, err)
	}

	return nil
}
