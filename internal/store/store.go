package store

import jwt "github.com/dgrijalva/jwt-go"

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserService interface {
	// User return user associated with the given ID
	User(id int64) (*User, error)

	// UserByName return user associated with given Username
	UserByName(username string) (*User, error)

	// UserByEmail return user associated with given Email
	UserByEmail(email string) (*User, error)

	// Add insert given user instance in database
	Add(user *User) error
}

type ValidationService interface {
	// Add create a new column for the given userid with validation code
	Add(userID int64, code string) error

	// IsValidated check if given user is valided
	IsValidated(userID int64) (bool, error)

	// CheckCode check if given code is valid
	CheckCode(userID int64, code string) (bool, error)

	// Validate consume user validation code
	Validate(userID int64, code string) error
}

type AuthTokenService interface {
	Add(tokenString string, token *jwt.Token)
	Delete(tokenString string)
	Exists(token string) bool
}

type Store struct {
	UserService       UserService
	ValidationService ValidationService
	AuthTokenService  AuthTokenService
}
