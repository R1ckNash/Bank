package auth

import (
	"auth/internal/app/models"
	"auth/internal/app/repository/user_storage"
	"context"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
	"github.com/R1ckNash/Bank/pkg/helpers"
	"github.com/R1ckNash/Bank/pkg/transaction_manager"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser - user registration
func (as *authService) RegisterUser(ctx context.Context, user *models.User) error {
	const api = "auth.RegisterUser"

	// TODO: idempotency

	log := as.Logger.With(zap.String("op", api), zap.String("email", user.Email))

	log.Info("registration new user")

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("hash password error", zap.Error(err))
		return pkgerrors.Wrap(api, err)
	}

	userDto := &user_storage.User{
		Email:    user.Email,
		Username: user.Username,
		Password: string(hashed),
	}

	err = helpers.WithRetries(ctx, func(ctx context.Context) error {
		var err error
		err = as.TransactionManager.RunReadCommitted(ctx, transaction_manager.ReadWrite,
			func(txCtx context.Context) error { // TRANSANCTION SCOPE
				if err = as.UserStorage.CreateUser(txCtx, userDto); err != nil {
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
