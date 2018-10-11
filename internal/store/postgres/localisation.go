package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type LocalisationService struct {
	db *sql.DB
}

const (
	createLocalisationTableSQL = `
		CREATE TABLE IF NOT EXISTS localisations (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			lat float NOT NULL,
			lng float NOT NULL,
			address varchar(256) NOT NULL
		);
	`
)

func NewLocalisationService(db *sql.DB) store.LocalisationService {
	if _, err := db.Exec(createLocalisationTableSQL); err != nil {
		panic(err)
	}
	return &LocalisationService{db}
}

func (s *LocalisationService) Add(l *store.Localisation) error {
	return s.db.QueryRow("INSERT INTO localisations (userId, lat, lng, address) VALUES($1,$2,$3,$4) RETURNING id",
		&l.UserID,
		&l.Lat,
		&l.Lng,
		&l.Address,
	).Scan(&l.ID)
}

func (s *LocalisationService) Localisation(userID int64) (*store.Localisation, error) {
	loc := &store.Localisation{}
	err := s.db.QueryRow("SELECT * FROM localisations WHERE userId=$1", userID).Scan(
		&loc.ID,
		&loc.UserID,
		&loc.Lat,
		&loc.Lng,
		&loc.Address,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return loc, err
}
