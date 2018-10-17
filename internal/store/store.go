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
	UserID     int64  `json:"userId"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Birthdate  int64  `json:"birthdate"`
	Bio        string `json:"bio"`
	Gender     string `json:"gender"`
	Attraction string `json:"attraction"`
}

type Picture struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"userId"`
	Path   string `json:"path"`
	IsPP   bool   `json:"isPP"`
}

type Localisation struct {
	ID      int64   `json:"-"`
	UserID  int64   `json:"userId"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Address string  `json:"address"`
}

type Notification struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type ValidationCode struct {
	Code string `json:"code"`
	Used bool   `json:"used"`
}

type HistoryItem struct {
	ID     int64  `json:"-"`
	UserID int64  `json:"userId"`
	FromID int64  `json:"fromId"`
	Event  string `json:"event"`
}

type Interest struct {
	ID     int64  `json:"-"`
	UserID int64  `json:"userId"`
	Value  string `json:"value"`
}

type Message struct {
	ID       int64  `json:"id"`
	Sender   int64  `json:"sender"`
	Reciever int64  `json:"reciever"`
	Message  string `json:"message"`
	Date     int64  `json:"date"`
}

type Like struct {
	ID       int64 `json:"-"`
	UserID   int64 `json:"-"`
	TargetID int64 `json:"-"`
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
	Profiles() ([]*Profile, error)
	Profile(userID int64) (*Profile, error)
	Add(profile *Profile) error
	Delete(userID int64) error
	Update(profile *Profile) error
}

type LikesService interface {
	Likes() ([]*Like, error)
	Like(userID, id int64) (bool, error)
	Add(userID, id int64) error
	Remove(userID, id int64) error
	Count(userID int64) (int, error)
}

type MatchService interface {
	Matches(userID int64) ([]int64, error)
	Match(usersID ...int64) (bool, error)
	Add(usersID ...int64) error
	Remove(usersID ...int64) error
}

type LocalisationService interface {
	Localisation(userID int64) (*Localisation, error)
	Localisations() ([]*Localisation, error)
	DistanceBetween(userID, targetID int64) (float64, error)
	Add(l *Localisation) error
	Update(l *Localisation) error
}

type InterestService interface {
	Add(interest *Interest) error
	Remove(userID int64, slug string) error
	Interest(userID int64, slug string) (*Interest, error)
	Interests() ([]*Interest, error)
	AllBySlug(slug string, limit int) ([]*Interest, error)
	AllByUser(userID int64) ([]*Interest, error)
}

type MessageService interface {
	Add(msg *Message) error
	Remove(id int64) error
	Messages(userId int64) ([]*Message, error)
	MessagesBetween(u1, u2 int64) ([]*Message, error)
}

type HistoryService interface {
	All() ([]*HistoryItem, error)
	Add(i *HistoryItem) error
	Exists(i *HistoryItem) (bool, error)
	Count(userID int64) (int, error)
	History(userID int64) ([]*HistoryItem, error)
}

type PresenceService interface {
	Add(userID int64) error
	Heartbeat(userID int64) error
	LastHeartbeat(userID int64) (int64, error)
}

type IgnoreService interface {
	Add(userID, target int64) error
	Remove(userID, target int64) error
	Ignored(userID, target int64) (bool, error)
}

type ReportService interface {
	Add(userID, target int64) error
	Reported(userID, target int64) (bool, error)
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

type NotificationService interface {
	Push(to int64, from int64, evt string, data interface{}) error
	Notifications(userID int64) []*Notification
}

type AuthTokenService interface {
	Token(tokenString string) (*jwt.Token, error)
	Add(tokenString string, token *jwt.Token) error
	Delete(tokenString string) error
}

type Store struct {
	UserService         UserService
	LocalisationService LocalisationService
	NotificationService NotificationService
	MessageService      MessageService
	MatchService        MatchService
	InterestService     InterestService
	LikesService        LikesService
	HistoryService      HistoryService
	PicturesService     PicturesService
	ValidationService   ValidationService
	IgnoreService       IgnoreService
	AuthTokenService    AuthTokenService
	ReportService       ReportService
	PresenceService     PresenceService
	ProfileService      ProfileService
}
