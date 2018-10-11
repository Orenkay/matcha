package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type InterestService struct {
	db *sql.DB
}

const (
	createInterestsTableSQL = `
		CREATE TABLE IF NOT EXISTS users_interests (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			value varchar(64) NOT NULL
		);
	`
)

func NewInterestService(db *sql.DB) store.InterestService {
	if _, err := db.Exec(createInterestsTableSQL); err != nil {
		panic(err)
	}
	return &InterestService{db}
}

func (s *InterestService) Add(interest *store.Interest) error {
	return s.db.QueryRow("INSERT INTO users_interests (userId, value) VALUES($1,$2) RETURNING id",
		&interest.UserID,
		&interest.Value,
	).Scan(&interest.ID)
}

func (s *InterestService) Remove(userID int64, slug string) error {
	_, err := s.db.Exec("DELETE FROM users_interests WHERE userId=$1 AND value=$2", userID, slug)
	return err
}

func (s *InterestService) AllBySlug(slug string, limit int) ([]*store.Interest, error) {
	var interests []*store.Interest
	rows, err := s.db.Query("SELECT value FROM users_interests WHERE value LIKE $1 GROUP BY value ORDER BY COUNT(*) DESC LIMIT $2", slug+"%", limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		interest := &store.Interest{}
		if err := rows.Scan(&interest.Value); err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}
	return interests, nil
}

func (s *InterestService) Interest(userID int64, slug string) (*store.Interest, error) {
	i := &store.Interest{}
	err := s.db.QueryRow("SELECT value FROM users_interests WHERE userId=$1 AND value=$2", userID, slug).Scan(&i.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return i, err
}

func (s *InterestService) AllByUser(userID int64) ([]*store.Interest, error) {
	var interests []*store.Interest
	rows, err := s.db.Query("SELECT value FROM users_interests WHERE userId=$1", userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		interest := &store.Interest{}
		if err := rows.Scan(&interest.Value); err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}
	return interests, nil
}
