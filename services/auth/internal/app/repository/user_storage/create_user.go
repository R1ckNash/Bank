package user_storage

import (
	"auth/internal/app/models"
	"context"
	"errors"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *UserStorage) CreateUser(ctx context.Context, user *User) error {
	const api = "user_storage.CreateUser"

	query := `insert into users (username, email, password) values ($1, $2, $3) returning id`

	if _, err := s.driver.GetQueryEngine(ctx).Exec(ctx, query, user.Username, user.Email, user.Password); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return pkgerrors.Wrap(api, models.ErrAlreadyExists)
		}
		return pkgerrors.Wrap(api, err)
	}

	return nil
}
