package main

import (
	"database/sql"

	"github.com/orenkay/matcha/internal/server"
	"github.com/orenkay/matcha/internal/store"
	"github.com/orenkay/matcha/internal/store/memory"
	"github.com/orenkay/matcha/internal/store/postgres"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=postgres user=matcha password=matcha dbname=matcha sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := &store.Store{
		UserService:       postgres.NewUserService(db),
		ValidationService: postgres.ValidationService(db),
		AuthToken:         memory.NewAuthTokenService(db),
	}

	server := server.New(store)
	server.Run("0.0.0.0:3000")
}
