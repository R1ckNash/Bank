package auth

import (
	"auth/domain"
	"context"
	"fmt"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
	"github.com/R1ckNash/Bank/pkg/helpers"
	"github.com/R1ckNash/Bank/pkg/transaction_manager"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//go:generate mockery --name=AuthService --filename=auth_service_mock.go --disable-version-string
type AuthService interface {
	RegisterUser(ctx context.Context, user *domain.User) error
	LoginUser(ctx context.Context, username, password string) (string, error)
	VerifyUser(ctx context.Context, id uuid.UUID) bool
}

//go:generate mockery --name=UserRepository --filename=user_repository_mock.go --disable-version-string
type UserRepository interface {
	StoreUser(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

//go:generate mockery --name=OutboxRepository --filename=outbox_repository_mock.go --disable-version-string
type OutboxRepository interface {
	AddMessage(ctx context.Context, msg *domain.OutboxMessage) error
	GetUnsentMessages(ctx context.Context, limit int) ([]*domain.OutboxMessage, error)
	MarkAsSent(ctx context.Context, id int64) error
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
	UserRepository
	OutboxRepository
	TransactionManager
	Producer  EventProducer
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

// RegisterUser - user registration
func (as *authService) RegisterUser(ctx context.Context, user *domain.User) error {
	const api = "auth.RegisterUser"

	// TODO: idempotency

	log := as.Logger.With(zap.String("op", api), zap.String("email", user.Email))

	log.Info("registration new user")

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("hash password error", zap.Error(err))
		return pkgerrors.Wrap(api, err)
	}

	user.Password = string(hashed)
	eventPayload := fmt.Sprintf(`{"event_type":"registration","user_id":"%s","email":"%s","status":"success"}`, user.ID, user.Email)

	err = helpers.WithRetries(ctx, func(ctx context.Context) error {
		var err error
		err = as.TransactionManager.RunReadCommitted(ctx, transaction_manager.ReadWrite,
			func(txCtx context.Context) error { // TRANSANCTION SCOPE
				if err = as.UserRepository.StoreUser(txCtx, user); err != nil {
					return err
				}

				outboxMsg := &domain.OutboxMessage{
					AggregateType: "user",
					AggregateID:   user.ID,
					Type:          "user_registered",
					Payload:       eventPayload,
					CreatedAt:     time.Now(),
				}
				if err := as.OutboxRepository.AddMessage(txCtx, outboxMsg); err != nil {
					return err
				}

				return nil
			},
		)
		return err
	})

	if err != nil {
		log.Warn("error saving user", zap.Error(err))
		return pkgerrors.Wrap(api, err)
	}

	return nil
}

// LoginUser - user log in
func (as *authService) LoginUser(ctx context.Context, username, password string) (string, error) {
	const api = "auth.LoginUser"

	log := as.Logger.With(zap.String("op", api), zap.String("username", username))

	log.Info("login user")

	user, err := as.UserRepository.GetByUsername(ctx, username)
	if err != nil {
		log.Error("user not found", zap.Error(err))
		return "", domain.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Error("user not found", zap.Error(err))
		return "", domain.ErrInvalidPassword
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(as.JwtSecret))
	if err != nil {
		log.Error("token generation error", zap.Error(err))
		return "", domain.ErrInvalidToken
	}

	return tokenString, nil
}

// VerifyUser - verification if user exists
func (as *authService) VerifyUser(ctx context.Context, id uuid.UUID) bool {
	const api = "auth.VerifyUser"

	log := as.Logger.With(zap.String("op", api), zap.String("id", id.String()))

	log.Info("verification user")

	_, err := as.UserRepository.GetByID(ctx, id)
	if err != nil {
		log.Error("user not found", zap.Error(err))
		return false
	}

	return true
}
