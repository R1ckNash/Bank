package models

import "github.com/google/uuid"

type Account struct {
	OwnerID  uuid.UUID `json:"owner_id"`
	Name     string    `json:"name"`
	Currency string    `json:"currency"`
	Email    string    `json:"email"`
}
