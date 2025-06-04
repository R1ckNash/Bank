package models

type Account struct {
	OwnerID  string `json:"owner_id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Email    string `json:"email"`
}
