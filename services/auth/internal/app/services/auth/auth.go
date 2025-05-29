package auth

import (
	"auth/internal/app/models"
	"auth/internal/app/repository/user_storage"
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

//go:generate mockery --name=AuthService --filename=auth_service_mock.go --disable-version-string
type AuthService interface {
	// RegisterUser - user creation
	RegisterUser(ctx context.Context, user *models.User) error
	// LoginUser - user log in
	LoginUser(ctx context.Context, username, password string) (string, error)
	// VerifyUser - if user exist
	VerifyUser(ctx context.Context, id int64) bool
}

//go:generate mockery --name=UserStorage --filename=user_storage_mock.go --disable-version-string
type UserStorage interface {
	// CreateUser - user creation
	CreateUser(ctx context.Context, user *user_storage.User) error
	// GetByUsername - get user by username
	GetByUsername(ctx context.Context, username string) (*user_storage.User, error)
	// GetByID - get user by id
	GetByID(ctx context.Context, userID int64) (*user_storage.User, error)
}

// TransactionManager trx manager
type TransactionManager interface {
	RunReadCommitted(ctx context.Context, accessMode pgx.TxAccessMode, f func(ctx context.Context) error) error
}

type Deps struct {
	UserStorage
	TransactionManager
	JwtSecret string
	Logger    *zap.Logger
}

type authService struct {
	Deps
}

func NewAuthService(d Deps) AuthService {
	return &authService{
		Deps: d,
	}
}
