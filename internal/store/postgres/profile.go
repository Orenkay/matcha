package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type ProfileService struct {
	db *sql.DB
}

const (
	createProfilesTableSQL = `
		CREATE TABLE IF NOT EXISTS profiles (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			lastName varchar(64) NOT NULL,
			firstName varchar(64) NOT NULL,
			gender varchar(32) NOT NULL,
			attraction varchar(32) NOT NULL,
			bio varchar(200) NOT NULL
		);
	`
)

func NewProfileService(db *sql.DB) store.ProfileService {
	if _, err := db.Exec(createProfilesTableSQL); err != nil {
		panic(err)
	}
	return &ProfileService{db}
}

func (s *ProfileService) Profile(userID int64) (*store.Profile, error) {
	p := &store.Profile{}
	err := s.db.QueryRow("SELECT * FROM profiles WHERE userId=$1", userID).Scan(
		&p.ID,
		&p.UserID,
		&p.LastName,
		&p.FirstName,
		&p.Gender,
		&p.Attraction,
		&p.Bio,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (s *ProfileService) Add(profile *store.Profile) error {
	return s.db.QueryRow("INSERT INTO profiles (userId, lastName, firstName, gender, attraction, bio) VALUES($1,$2,$3,$4,$5,$6) RETURNING id",
		profile.UserID,
		profile.LastName,
		profile.FirstName,
		profile.Gender,
		profile.Attraction,
		profile.Bio,
	).Scan(&profile.ID)
}

func (s *ProfileService) Delete(userId int64) error {
	return s.db.QueryRow("DELETE FROM profiles WHERE userId=$1", userId).Scan()
}

func (s *ProfileService) Update(profile *store.Profile) error {
	_, err := s.db.Exec("UPDATE profiles SET lastName=$1, firstName=$2, gender=$3, attraction=$4, bio=$5 WHERE userId=$6",
		profile.LastName,
		profile.FirstName,
		profile.Gender,
		profile.Attraction,
		profile.Bio,
		profile.UserID,
	)
	return err
}
