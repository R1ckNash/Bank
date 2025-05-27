package user_storage

import (
	"Bank/pkg/transaction_manager"
	"context"
)

type UserStorage struct {
	// connection *postgres.Connection // если тетсируте только интеграционными
	// connection Connection // если мокаете базу данных
	driver QueryEngineProvider
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) transaction_manager.QueryEngine
}

// New - returns UserStorage
func New( /*connection *postgres.Connection*/ driver QueryEngineProvider) *UserStorage {
	return &UserStorage{
		// connection: connection,
		driver: driver,
	}
}
