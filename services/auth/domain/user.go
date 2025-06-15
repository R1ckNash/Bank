package domain

import (
	"github.com/google/uuid"
	"time"
)

// User is representing the User data struct
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
