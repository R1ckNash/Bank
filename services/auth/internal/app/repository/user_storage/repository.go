package user_storage

import (
	"context"
	"github.com/R1ckNash/Bank/pkg/transaction_manager"
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
