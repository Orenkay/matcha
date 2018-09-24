package store

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

	// IsUsed check if given user is valided
	IsUsed(userID int64) (bool, error)

	// Code return validation code associated with the given UserID
	Code(userID int64) (string, error)

	// Use consume user validation code
	Valid(userID int64) error
}

type AuthTokenService interface {
	Add(userID int64, token string) error
	Delete(token string) error
	Exists(token string) (bool, error)
}

type Store struct {
	UserService       UserService
	ValidationService ValidationService
	AuthTokenService  AuthTokenService
}
