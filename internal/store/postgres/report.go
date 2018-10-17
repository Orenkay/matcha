package postgres

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/store"
)

type ReportService struct {
	db *sql.DB
}

const (
	createReportTableSQL = `
		CREATE TABLE IF NOT EXISTS users_reports (
			id SERIAL PRIMARY KEY,
			userId int NOT NULL,
			targetId int NOT NULL
		);
	`
)

func NewReportService(db *sql.DB) store.ReportService {
	if _, err := db.Exec(createReportTableSQL); err != nil {
		panic(err)
	}
	return &ReportService{db}
}

func (s *ReportService) Add(userID, target int64) error {
	_, err := s.db.Exec("INSERT INTO users_reports(userId, targetId) VALUES($1, $2)", &userID, &target)
	return err
}

func (s *ReportService) Reported(userID, target int64) (bool, error) {
	err := s.db.QueryRow("SELECT FROM users_reports WHERE userId=$1 AND targetId=$2", &userID, &target).Scan()
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
