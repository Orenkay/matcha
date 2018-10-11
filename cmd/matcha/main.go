package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/orenkay/matcha/internal/server"
	"github.com/orenkay/matcha/internal/store"
	"github.com/orenkay/matcha/internal/store/memory"
	"github.com/orenkay/matcha/internal/store/postgres"
)

func main() {
	db, err := sql.Open("postgres", "host=postgres user=matcha password=matcha dbname=matcha sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := &store.Store{
		UserService:         postgres.NewUserService(db),
		ValidationService:   postgres.NewValidationService(db),
		InterestService:     postgres.NewInterestService(db),
		LocalisationService: postgres.NewLocalisationService(db),
		ProfileService:      postgres.NewProfileService(db),
		AuthTokenService:    memory.NewAuthTokenService(),
	}

	server := server.New(store)
	server.Run("0.0.0.0:3000")
}
