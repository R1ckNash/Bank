package user_storage

import (
	"context"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
)

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (*User, error) {
	const api = "user_storage.GetByUsername"

	query := `select id, username, email, password from users where username=$1`
	row := s.driver.GetQueryEngine(ctx).QueryRow(ctx, query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, pkgerrors.Wrap(api, err)
	}

	return &user, nil
}
