package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type HistoryService struct {
	db *sql.DB
}

const (
	createHistoryTableSQL = `
		CREATE TABLE IF NOT EXISTS users_history (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			fromId int NOT NULL,
			event varchar(32) NOT NULL
		);
	`
)

func NewHistoryService(db *sql.DB) store.HistoryService {
	if _, err := db.Exec(createHistoryTableSQL); err != nil {
		panic(err)
	}
	return &HistoryService{db}
}

func (s *HistoryService) All() ([]*store.HistoryItem, error) {
	var items []*store.HistoryItem

	rows, err := s.db.Query("SELECT * FROM users_history")
	{
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		i := &store.HistoryItem{}
		if err := rows.Scan(&i.ID, &i.UserID, &i.FromID, &i.Event); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (s *HistoryService) Add(i *store.HistoryItem) error {
	_, err := s.db.Exec("INSERT INTO users_history(userId, fromId, event) VALUES($1, $2, $3)", &i.UserID, &i.FromID, &i.Event)
	return err
}

func (s *HistoryService) Count(userID int64) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) WHERE users_history WHERE userId=$1", userID).Scan(&count)
	{
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, nil
			}
			return 0, err
		}
	}
	return count, nil
}

func (s *HistoryService) Exists(i *store.HistoryItem) (bool, error) {
	err := s.db.QueryRow("SELECT FROM users_history WHERE (userId=$1 AND fromId=$2 AND event=$3)", &i.UserID, &i.FromID, &i.Event).Scan()
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

func (s *HistoryService) History(userID int64) ([]*store.HistoryItem, error) {
	var history []*store.HistoryItem
	rows, err := s.db.Query("SELECT * FROM users_history WHERE userId=$1 AND profile_check(fromId)", &userID)
	{
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		i := &store.HistoryItem{}
		if err := rows.Scan(&i.ID, &i.UserID, &i.FromID, &i.Event); err != nil {
			return nil, err
		}
		history = append(history, i)
	}
	return history, nil
}
