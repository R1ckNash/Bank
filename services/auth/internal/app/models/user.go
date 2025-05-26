package models

type User struct {
	ID           UserID
	Username     string
	PasswordHash string
	Email        string
}
