package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type LocalisationService struct {
	db *sql.DB
}

const (
	// http://en.wikipedia.org/wiki/Haversine_formula
	createCoordDistanceFuncSQL = `
		DROP FUNCTION IF EXISTS coord_distance;
		CREATE OR REPLACE FUNCTION coord_distance(olat float, dlat float, olng float, dlng float) RETURNS float AS $$
		BEGIN
			RETURN (
			2 * 6378100 * ASIN
				(
					SQRT
					(
						POWER
						(
							SIN((olat - dlat) * pi() / 180 / 2),
							2
						)
						+
						COS(olat * pi() / 180)
						*
						COS(dlat * pi() / 180)
						*
						POWER
						(
							SIN((olng - dlng) * pi() / 180 / 2),
							2
						)
					)
				)
			);
		END;
		$$ LANGUAGE plpgsql;
	`
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
	if _, err := db.Exec(createCoordDistanceFuncSQL); err != nil {
		panic(err)
	}
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

func (s *LocalisationService) Update(l *store.Localisation) error {
	_, err := s.db.Exec("UPDATE localisations SET lat=$1, lng=$2, address=$3 WHERE userId=$4",
		&l.Lat,
		&l.Lng,
		&l.Address,
		&l.UserID,
	)
	return err
}

func (s *LocalisationService) DistanceBetween(userID, targetID int64) (float64, error) {
	var distance float64
	err := s.db.QueryRow(`
		SELECT
			COORD_DISTANCE(t1.lat, t2.lat, t1.lng, t2.lng) as distance
		FROM localisations t1
		INNER JOIN localisations t2 ON t2.userId=$2
		WHERE t1.userId=$1;
	`, userID, targetID).Scan(&distance)
	return distance, err
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

func (s *LocalisationService) Localisations() ([]*store.Localisation, error) {
	var locs []*store.Localisation
	rows, err := s.db.Query("SELECT * FROM localisations")
	{
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		l := &store.Localisation{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.Lat, &l.Lng, &l.Address); err != nil {
			return nil, err
		}
		locs = append(locs, l)
	}
	return locs, nil
}
