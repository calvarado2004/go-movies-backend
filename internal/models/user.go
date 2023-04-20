package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is a struct that holds the user information.
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// PasswordMatches compares the plain text password with the hashed password.
func (u *User) PasswordMatches(plainText string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// Passwords don't match
			return false, nil
		default:
			return false, err

		}
	}

	return true, nil
}
