package user_storage

import (
	pkgerrors "Bank/pkg/errors"
	"auth/internal/app/models"
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *UserStorage) CreateUser(ctx context.Context, user *User) error {
	const api = "user_storage.CreateUser"

	query := `insert into users (email, password) values ($1, $2) returning id`

	if _, err := r.driver.GetQueryEngine(ctx).Exec(ctx, query, user.Email, user.Password); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return pkgerrors.Wrap(api, models.ErrAlreadyExists)
		}
		return pkgerrors.Wrap(api, err)
	}

	return nil
}
