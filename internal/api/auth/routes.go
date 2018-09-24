package auth

import (
	"github.com/go-chi/chi"
	"github.com/orenkay/matcha/internal/store"
)

func Routes(store *store.Store) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/login", Authenticate(store))
	router.Post("/register", Register(store))
	return router
}
