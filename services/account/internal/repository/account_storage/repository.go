package account_storage

import (
	"context"
	"github.com/R1ckNash/Bank/pkg/transaction_manager"
)

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) transaction_manager.QueryEngine
}

type AccountStorage struct {
	driver QueryEngineProvider
}

func New(driver QueryEngineProvider) *AccountStorage {
	return &AccountStorage{driver: driver}
}
