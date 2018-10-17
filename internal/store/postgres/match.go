package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type MatchService struct {
	db *sql.DB
	ns store.NotificationService
}

const (
	createMatchTableSQL = `
		CREATE TABLE IF NOT EXISTS users_matches (
			id SERIAL PRIMARY KEY,
			user1Id int NOT NULL,
			user2Id int NOT NULL
		);
	`
)

func NewMatchService(db *sql.DB, ns store.NotificationService) store.MatchService {
	if _, err := db.Exec(createMatchTableSQL); err != nil {
		panic(err)
	}
	return &MatchService{
		db: db,
		ns: ns,
	}
}

func (s *MatchService) Matches(userID int64) ([]int64, error) {
	var matches []int64

	rows, err := s.db.Query("SELECT user1Id, user2Id FROM users_matches WHERE user1Id=$1 OR user2Id=$1", userID)
	{
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var user1ID, user2ID int64
		if err := rows.Scan(&user1ID, &user2ID); err != nil {
			return nil, err
		}
		if user1ID == userID {
			matches = append(matches, user2ID)
		}
		if user2ID == userID {
			matches = append(matches, user1ID)
		}
	}
	return matches, nil
}

func (s *MatchService) Match(usersID ...int64) (bool, error) {
	if err := s.db.QueryRow("SELECT FROM users_matches WHERE (user1Id=$1 AND user2Id=$2) OR (user1Id=$2 AND user2Id=$1)", usersID[0], usersID[1]).Scan(); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *MatchService) Add(usersID ...int64) error {
	_, err := s.db.Exec("INSERT INTO users_matches(user1Id, user2Id) VALUES($1, $2)", usersID[0], usersID[1])
	for i, id := range usersID {
		s.ns.Push(id, usersID[1-i], "match", usersID[1-i])
	}
	return err
}

func (s *MatchService) Remove(usersID ...int64) error {
	res, err := s.db.Exec("DELETE FROM users_matches WHERE (user1Id=$1 AND user2Id=$2) OR (user1Id=$2 AND user2Id=$1)", usersID[0], usersID[1])
	{
		if err != nil {
			return err
		}
	}
	c, _ := res.RowsAffected()
	if c > 0 {
		for i, id := range usersID {
			s.ns.Push(id, usersID[1-i], "unmatch", usersID[1-i])
		}
	}
	return err
}
