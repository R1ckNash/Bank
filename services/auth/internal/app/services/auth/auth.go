package auth

import (
	"auth/internal/app/models"
	"auth/internal/app/repository/user_storage"
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type AuthService interface {
	// RegisterUser - user creation
	RegisterUser(ctx context.Context, user *models.User) error
}

//go:generate mockery --name=UserStorage --filename=user_storage_mock.go --disable-version-string
type UserStorage interface {
	// CreateUser - user creation
	//
	// @errors: models.ErrAlreadyExists
	CreateUser(ctx context.Context, user *user_storage.User) error
}

// TransactionManager trx manager
type TransactionManager interface {
	RunReadCommitted(ctx context.Context, accessMode pgx.TxAccessMode, f func(ctx context.Context) error) error
}

type Deps struct {
	UserStorage
	TransactionManager
	logger *zap.Logger
}

type authService struct {
	Deps
}

func NewAuthService(d Deps) AuthService {
	return &authService{
		Deps: d,
	}
}
