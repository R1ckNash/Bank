package postgres

import (
	"auth/domain"
	"context"
	"errors"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
	"github.com/R1ckNash/Bank/pkg/transaction_manager"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserStorage struct {
	driver QueryEngineProvider
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) transaction_manager.QueryEngine
}

// New - returns UserStorage
func New(driver QueryEngineProvider) *UserStorage {
	return &UserStorage{
		driver: driver,
	}
}

func (s *UserStorage) fetch(ctx context.Context, query string, args ...interface{}) (result []*domain.User, err error) {
	rows, err := s.driver.GetQueryEngine(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result = make([]*domain.User, 0)

	for rows.Next() {
		t := domain.User{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.Email,
			&t.Password,
			&t.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, &t)
	}
	return result, nil
}

func (s *UserStorage) StoreUser(ctx context.Context, user *domain.User) error {
	const api = "postgres.StoreUser"

	query := `insert into users (id, username, email, password, created_at) values ($1, $2, $3, $4, $5) returning id`

	if _, err := s.driver.GetQueryEngine(ctx).Exec(ctx, query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return pkgerrors.Wrap(api, domain.ErrAlreadyExists)
		}
		return pkgerrors.Wrap(api, err)
	}

	return nil
}

func (s *UserStorage) GetByID(ctx context.Context, userID uuid.UUID) (res *domain.User, err error) {
	const api = "postgres.GetByID"
	query := `SELECT id, username, email, password, created_at FROM users WHERE id=$1`
	list, err := s.fetch(ctx, query, userID)
	if err != nil {
		return nil, pkgerrors.Wrap(api, err)
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, pkgerrors.Wrap(api, domain.ErrUserNotFound)
	}

	return
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (res *domain.User, err error) {
	const api = "postgres.GetByUsername"

	query := `select id, username, email, password, created_at from users where username=$1`
	list, err := s.fetch(ctx, query, username)
	if err != nil {
		return nil, pkgerrors.Wrap(api, err)
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, pkgerrors.Wrap(api, domain.ErrUserNotFound)
	}

	return
}
