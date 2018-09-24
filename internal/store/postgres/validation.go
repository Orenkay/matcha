package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type ValidationService struct {
	db *sql.DB
}

const (
	createvalidationTableSQL = `
		CREATE TABLE IF NOT EXISTS validations (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			code varchar(64) NOT NULL,
			used boolean DEFAULT FALSE
		);
	`
)

func NewValidationService(db *sql.DB) store.ValidationService {
	if _, err := db.Exec(createvalidationTableSQL); err != nil {
		panic(err)
	}
	return &ValidationService{db}
}

// Add create a new column for the given userid with validation code
func (s *ValidationService) Add(userID int64, code string) error {
	_, err := s.db.Exec("INSERT INTO validations (userId, code) VALUES($1, $2)", userID, code)
	return err
}

// IsUsed check if given user is valided
func (s *ValidationService) IsUsed(userID int64) (bool, error) {
	var valid bool
	err := s.db.QueryRow("SELECT used FROM validations WHERE userId=$1", userID).Scan(&valid)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Code return validation code associated with the given UserID
func (s *ValidationService) Code(userID int64) (string, error) {
	var userCode string

	err := s.db.QueryRow("SELECT code FROM validations WHERE userId=$1", userID).Scan(&userCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return userCode, nil
}

// Valid consume user validation code
func (s *ValidationService) Valid(userID int64) error {
	_, err := s.db.Exec("UPDATE validations SET used=true WHERE userId=$1", userID)
	return err
}
