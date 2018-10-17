package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type IgnoreService struct {
	db *sql.DB
}

const (
	createIgnoresTableSQL = `
		CREATE TABLE IF NOT EXISTS users_ignores (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			targetId int NOT NULL
		);
	`
)

func NewIgnoreService(db *sql.DB) store.IgnoreService {
	if _, err := db.Exec(createIgnoresTableSQL); err != nil {
		panic(err)
	}
	return &IgnoreService{db}
}

func (s *IgnoreService) Add(userID, target int64) error {
	_, err := s.db.Exec("INSERT INTO users_ignores(userId, targetId) VALUES($1, $2)", userID, target)
	return err
}

func (s *IgnoreService) Remove(userID, target int64) error {
	_, err := s.db.Exec("DELETE FROM users_ignores WHERE userId=$1 AND targetId=$2", userID, target)
	return err
}

func (s *IgnoreService) Ignored(userID, target int64) (bool, error) {
	err := s.db.QueryRow("SELECT FROM users_ignores WHERE userId=$1 AND targetId=$2", userID, target).Scan()
	{
		if err != nil {
			if err == sql.ErrNoRows {
				return false, nil
			}
			return false, err
		}
	}
	return true, nil
}
