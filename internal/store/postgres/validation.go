package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type ValidationService struct {
	db *sql.DB
}

const (
	createValidationTableSQL = `
		CREATE TABLE IF NOT EXISTS validations (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			code varchar(64) NOT NULL,
			used boolean DEFAULT FALSE
		);
	`
)

func NewValidationService(db *sql.DB) store.ValidationService {
	if _, err := db.Exec(createValidationTableSQL); err != nil {
		panic(err)
	}
	return &ValidationService{db}
}

// Add create a new column for the given userid with validation code
func (s *ValidationService) Add(userID int64, code string) error {
	_, err := s.db.Exec("INSERT INTO validations (userId, code) VALUES($1, $2)", userID, code)
	return err
}

func (s *ValidationService) Remove(userID int64) error {
	_, err := s.db.Exec("DELETE FROM validations WHERE userId=$1", userID)
	return err
}

// IsValidated check if given user is valided
func (s *ValidationService) IsValidated(userID int64) (bool, error) {
	var valid bool
	err := s.db.QueryRow("SELECT used FROM validations WHERE userId=$1", userID).Scan(&valid)
	if err != nil {
		return false, err
	}
	return valid, nil
}

func (s *ValidationService) CheckCode(userID int64, code string) (bool, error) {
	var code2 string

	err := s.db.QueryRow("SELECT code FROM validations WHERE userId=$1", userID).Scan(&code2)
	if err != nil {
		return false, err
	}
	return code == code2, nil
}

func (s *ValidationService) ValidationCode(userID int64) (*store.ValidationCode, error) {
	v := &store.ValidationCode{}

	err := s.db.QueryRow("SELECT code, used FROM validations WHERE userId=$1", userID).Scan(&v.Code, &v.Used)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Validate consume user validation code
func (s *ValidationService) Validate(userID int64, code string) error {
	_, err := s.db.Exec("UPDATE validations SET used=true WHERE userId=$1 AND code=$2", userID, code)
	return err
}
