package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type UserService struct {
	db *sql.DB
}

const (
	createUsersTableSQL = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email varchar(64) NOT NULL,
			username varchar(32) NOT NULL,
			password varchar(60) NOT NULL
		);
	`
)

func NewUserService(db *sql.DB) store.UserService {
	if _, err := db.Exec(createUsersTableSQL); err != nil {
		panic(err)
	}
	return &UserService{db}
}

// User return user associated with the given ID
func (s *UserService) User(id int64) (*store.User, error) {
	user := &store.User{ID: id}
	err := s.db.QueryRow("SELECT email, username, password FROM users WHERE id=$1", id).Scan(
		&user.Email, &user.Username, &user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// UserByName return user associated with given Username
func (s *UserService) UserByName(username string) (*store.User, error) {
	user := &store.User{Username: username}
	err := s.db.QueryRow("SELECT id, email, password FROM users WHERE username=$1", username).Scan(
		&user.ID, &user.Email, &user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// UserByEmail return user associated with given Email
func (s *UserService) UserByEmail(email string) (*store.User, error) {
	user := &store.User{Email: email}
	err := s.db.QueryRow("SELECT id, username, password FROM users WHERE email=$1", email).Scan(
		&user.ID, &user.Username, &user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Add insert given user instance in database
func (s *UserService) Add(user *store.User) error {
	r, err := s.db.Exec("INSERT INTO users (email, username, password) VALUES($1,$2,$3)", user.Email, user.Username, user.Password)
	if err != nil {
		return err
	}
	user.ID, _ = r.LastInsertId()
	return nil
}
