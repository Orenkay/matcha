package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type LikesService struct {
	db *sql.DB
	ns store.NotificationService
}

const (
	createLikesTableSQL = `
		CREATE TABLE IF NOT EXISTS users_likes (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			targetId int NOT NULL
		);
	`
)

func NewLikesService(db *sql.DB, ns store.NotificationService) store.LikesService {
	if _, err := db.Exec(createLikesTableSQL); err != nil {
		panic(err)
	}
	return &LikesService{db, ns}
}

func (s *LikesService) Likes() ([]*store.Like, error) {
	var likes []*store.Like
	rows, err := s.db.Query("SELECT * FROM users_likes")
	{
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		l := &store.Like{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.TargetID); err != nil {
			return nil, err
		}
		likes = append(likes, l)
	}
	return likes, nil
}

func (s *LikesService) Like(userID, id int64) (bool, error) {
	if err := s.db.QueryRow("SELECT FROM users_likes WHERE userId=$1 AND targetId=$2", userID, id).Scan(); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *LikesService) Add(userID, id int64) error {
	_, err := s.db.Exec("INSERT INTO users_likes(userId, targetId) VALUES($1, $2)", userID, id)
	s.ns.Push(id, userID, "like", userID)
	return err
}

func (s *LikesService) Count(userID int64) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users_likes WHERE userId=$1", userID).Scan(&count)
	{
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, nil
			}
			return 0, err
		}
	}
	return count, err
}

func (s *LikesService) Remove(userID, id int64) error {
	_, err := s.db.Exec("DELETE FROM users_likes WHERE userId=$1 AND targetId=$2", userID, id)
	s.ns.Push(id, userID, "unlike", userID)
	return err
}
