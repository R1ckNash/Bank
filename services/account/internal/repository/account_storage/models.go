package account_storage

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID        int64     `db:"id"`
	OwnerID   uuid.UUID `db:"owner_id"`
	Name      string    `db:"name"`
	Currency  string    `db:"currency"`
	Email     string    `db:"email"`
	IsBlocked bool      `db:"is_blocked"`
	Balance   float64   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
