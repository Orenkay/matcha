package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type UserService struct {
	db *sql.DB
}

const (
	profileCheckFuncSQL = `
		DROP FUNCTION IF EXISTS profile_check;
		CREATE OR REPLACE FUNCTION profile_check(_id int) RETURNS bool AS $$
		DECLARE passed BOOLEAN;
		BEGIN
			SELECT t1.userId = _id FROM profiles t1
				LEFT JOIN pictures c1 ON c1.userId = t1.userId AND isPP=true
				LEFT JOIN users_interests c2 ON c2.userId = t1.userId
				WHERE t1.userId = _id
				GROUP BY t1.userId, c1.userId, c2.userId
				HAVING (COUNT(c1) > 0 AND COUNT(c2) > 0) INTO passed;
			RETURN passed;
		END;
		$$ LANGUAGE plpgsql;
	`

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
	if _, err := db.Exec(profileCheckFuncSQL); err != nil {
		panic(err)
	}
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
	err := s.db.QueryRow("INSERT INTO users (email, username, password) VALUES($1,$2,$3) RETURNING id", user.Email, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Remove(user *store.User) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id=$1", user.ID)
	return err
}

func (s *UserService) Update(user *store.User) error {
	_, err := s.db.Exec("UPDATE users SET email=$1, username=$2, password=$3 WHERE id=$4",
		user.Email,
		user.Username,
		user.Password,
		user.ID,
	)
	return err
}
