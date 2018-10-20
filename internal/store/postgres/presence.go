package postgres

import (
	"database/sql"
	"time"

	"github.com/orenkay/matcha/internal/store"
)

type PresenceService struct {
	db *sql.DB
}

const (
	createPresenceTableSQL = `
		CREATE TABLE IF NOT EXISTS users_presences (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			time int NOT NULL
		);
	`
)

func NewPresenceService(db *sql.DB) store.PresenceService {
	if _, err := db.Exec(createPresenceTableSQL); err != nil {
		panic(err)
	}
	return &PresenceService{db}
}

func (s *PresenceService) Add(userID int64) error {
	_, err := s.db.Exec("INSERT INTO users_presences(userId, time) VALUES($1, $2)", &userID, time.Now().Unix())
	return err
}

func (s *PresenceService) Heartbeat(userID int64) error {
	_, err := s.db.Exec("UPDATE users_presences SET time=$1 WHERE userId=$2", time.Now().Unix(), &userID)
	{
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PresenceService) LastHeartbeat(userID int64) (int64, error) {
	var last int64
	err := s.db.QueryRow("SELECT time FROM users_presences WHERE userId=$1", &userID).Scan(&last)
	return last, err
}
