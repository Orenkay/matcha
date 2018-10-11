package store

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Profile struct {
	ID         int64  `json:"-"`
	UserID     int64  `json:"-"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Bio        string `json:"bio"`
	Gender     string `json:"gender"`
	Attraction string `json:"attraction"`
}

type Picture struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"-"`
	Path   string `json:"path"`
	IsPP   bool   `json:"isPP"`
}

type Localisation struct {
	ID      int64   `json:"-"`
	UserID  int64   `json:"-"`
	Lat     float64 `json:"-"`
	Lng     float64 `json:"-"`
	Address string  `json:"address"`
}

type ValidationCode struct {
	Code string `json:"code"`
	Used bool   `json:"used"`
}

type Interest struct {
	ID     int64  `json:"-"`
	UserID int64  `json:"-"`
	Value  string `json:"value"`
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

	Update(user *User) error
	Remove(user *User) error
}

type PicturesService interface {
	Add(p *Picture) error
	Update(p *Picture) error
	Delete(userID, id int64) error
	DeleteByPath(path string) error
	Pictures(userID int64) ([]*Picture, error)
	PicturesCount(userID int64) (int, error)
	Picture(userID, id int64) (*Picture, error)
	PP(userID int64) (*Picture, error)
}

type ProfileService interface {
	Profile(userID int64) (*Profile, error)
	Add(profile *Profile) error
	Delete(userID int64) error
	Update(profile *Profile) error
}

type LocalisationService interface {
	Localisation(userID int64) (*Localisation, error)
	Add(l *Localisation) error
	Update(l *Localisation) error
}

type InterestService interface {
	Add(interest *Interest) error
	Remove(userID int64, slug string) error
	Interest(userID int64, slug string) (*Interest, error)
	AllBySlug(slug string, limit int) ([]*Interest, error)
	AllByUser(userID int64) ([]*Interest, error)
}

type ValidationService interface {
	// Add create a new column for the given userid with validation code
	Add(userID int64, code string) error

	// IsValidated check if given user is valided
	IsValidated(userID int64) (bool, error)

	// ValidationCode return the validation column for the given user
	ValidationCode(userID int64) (*ValidationCode, error)

	// CheckCode check if given code is valid
	CheckCode(userID int64, code string) (bool, error)

	// Validate consume user validation code
	Validate(userID int64, code string) error

	// Remove user validation code in database
	Remove(userID int64) error
}

type AuthTokenService interface {
	Add(tokenString string, token *jwt.Token)
	Delete(tokenString string)
	Exists(token string) bool
}

type Store struct {
	UserService         UserService
	LocalisationService LocalisationService
	InterestService     InterestService
	PicturesService     PicturesService
	ValidationService   ValidationService
	AuthTokenService    AuthTokenService
	ProfileService      ProfileService
}
