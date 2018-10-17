package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type PicturesService struct {
	db *sql.DB
}

const (
	createPicturesTableSQL = `
		CREATE TABLE IF NOT EXISTS pictures (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			path varchar(256) NOT NULL,
			isPP boolean NOT NULL
		);
	`
)

func NewPicturesService(db *sql.DB) store.PicturesService {
	if _, err := db.Exec(createPicturesTableSQL); err != nil {
		panic(err)
	}
	return &PicturesService{db}
}

func (s *PicturesService) Add(p *store.Picture) error {
	return s.db.QueryRow("INSERT INTO pictures (userId, path, isPP) VALUES($1,$2,$3) RETURNING ID", p.UserID, p.Path, p.IsPP).Scan(&p.ID)
}

func (s *PicturesService) Update(p *store.Picture) error {
	_, err := s.db.Exec("UPDATE pictures SET path=$1, isPP=$2 WHERE id=$3", p.Path, p.IsPP, p.ID)
	return err
}

func (s *PicturesService) Picture(userID, id int64) (*store.Picture, error) {
	p := &store.Picture{}
	err := s.db.QueryRow("SELECT * FROM pictures WHERE id=$1 AND userId=$2", id, userID).Scan(
		&p.ID,
		&p.UserID,
		&p.Path,
		&p.IsPP,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (s *PicturesService) PP(userID int64) (*store.Picture, error) {
	p := &store.Picture{}
	err := s.db.QueryRow("SELECT * FROM pictures WHERE userId=$1 AND isPP=$2", userID, true).Scan(
		&p.ID,
		&p.UserID,
		&p.Path,
		&p.IsPP,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (s *PicturesService) Delete(userID, id int64) error {
	_, err := s.db.Exec("DELETE FROM pictures WHERE id=$1 AND userId=$2", id, userID)
	return err
}

func (s *PicturesService) DeleteByPath(path string) error {
	_, err := s.db.Exec("DELETE FROM pictures WHERE path=$1", path)
	return err
}

func (s *PicturesService) PicturesCount(userID int64) (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) as count FROM pictures WHERE userId=$1", userID).Scan(&count)
	{
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (s *PicturesService) Pictures(userID int64) ([]*store.Picture, error) {
	var pictures []*store.Picture
	rows, err := s.db.Query("SELECT * FROM pictures WHERE userId=$1 ORDER BY id ASC", userID)
	{
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		p := &store.Picture{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Path, &p.IsPP)
		{
			if err != nil {
				return nil, err
			}
		}
		pictures = append(pictures, p)
	}
	return pictures, nil
}
