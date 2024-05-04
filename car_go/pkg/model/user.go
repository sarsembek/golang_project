package model

import (
	"golang.org/x/crypto/bcrypt"
)

// User struct to represent a user in the system
type User struct {
	ID       int    `json:"id"`       // Unique identifier for the user
	Username string `json:"username"` // User's username, must be unique
	Password string `json:"password"` // User's hashed password
	// Additional fields like Email, Roles can be added here
}

// CreateUser hashes the password and creates a new user instance
func (u *User) CreateUser(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Authenticate checks if the provided password is correct
func (u *User) Authenticate(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// If there is an error (passwords don't match), return false and the error
		return false, err
	}
	// If no error, authentication is successful
	return true, nil
}
