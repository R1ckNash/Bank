package auth

import (
	"context"
	"go.uber.org/zap"
	"strconv"
)

// VerifyUser - verification if user exists
func (as *authService) VerifyUser(ctx context.Context, id int64) bool {
	const api = "auth.VerifyUser"

	log := as.Logger.With(zap.String("op", api), zap.String("id", strconv.FormatInt(id, 10)))

	log.Info("verification user")

	_, err := as.UserStorage.GetByID(ctx, id)
	if err != nil {
		log.Error("user not found", zap.Error(err))
		return false
	}

	return true
}
