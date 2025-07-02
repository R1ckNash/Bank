package account

import (
	"account/internal/models"
	"account/internal/repository/account_storage"
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

//go:generate mockery --name=AccountService --filename=account_service_mock.go --disable-version-string
type AccountService interface {
	RegisterAccount(ctx context.Context, account *models.Account) error
	GetAccount(ctx context.Context, accountID int64) (models.Account, error)
}

//go:generate mockery --name=AccountStorage --filename=account_storage_mock.go --disable-version-string
type AccountStorage interface {
	CreateAccount(ctx context.Context, account *account_storage.Account) error
	GetByID(ctx context.Context, accountID int64) (account_storage.Account, error)
}

//go:generate mockery --name=EventProducer
type EventProducer interface {
	SendMessage(topic, key string, message []byte) error
}

// TransactionManager trx manager
type TransactionManager interface {
	RunReadCommitted(ctx context.Context, accessMode pgx.TxAccessMode, f func(ctx context.Context) error) error
}

type Deps struct {
	AccountStorage
	TransactionManager
	EventProducer
	Logger *slog.Logger
}

type accountService struct {
	Deps
}

func NewAccountService(d Deps) AccountService {
	return &accountService{
		Deps: d,
	}
}
