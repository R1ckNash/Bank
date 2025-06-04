package account_storage

import (
	"context"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
)

func (s *AccountStorage) GetByID(ctx context.Context, accountID int64) (Account, error) {
	const api = "account_storage.GetByID"
	query := `SELECT name, owner_id, currency, email FROM account WHERE id=$1`
	row := s.driver.GetQueryEngine(ctx).QueryRow(ctx, query, accountID)

	var account Account
	err := row.Scan(&account.Name, &account.OwnerID, &account.Currency, &account.Email)
	if err != nil {
		return Account{}, pkgerrors.Wrap(api, err)
	}

	return account, nil
}
