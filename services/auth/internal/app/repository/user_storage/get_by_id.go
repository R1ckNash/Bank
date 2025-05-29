package user_storage

import (
	"context"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
)

func (s *UserStorage) GetByID(ctx context.Context, userID int64) (*User, error) {
	const api = "user_storage.GetByID"
	query := `SELECT id, username, email FROM users WHERE id=$1`
	row := s.driver.GetQueryEngine(ctx).QueryRow(ctx, query, userID)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, pkgerrors.Wrap(api, err)
	}

	return &user, nil
}
