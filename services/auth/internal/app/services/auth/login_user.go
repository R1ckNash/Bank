package auth

import (
	"auth/internal/app/models"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// LoginUser - user log in
func (as *authService) LoginUser(ctx context.Context, username, password string) (string, error) {
	const api = "auth.LoginUser"

	log := as.Logger.With(zap.String("op", api), zap.String("username", username))

	log.Info("login user")

	user, err := as.UserStorage.GetByUsername(ctx, username)
	if err != nil {
		log.Error("user not found", zap.Error(err))
		return "", models.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Error("user not found", zap.Error(err))
		return "", models.ErrInvalidPassword
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
		return "", models.ErrInvalidToken
	}

	return tokenString, nil
}
