package auth

import (
	"github.com/go-chi/chi"
	"github.com/orenkay/matcha/internal/api/middlewares"
	"github.com/orenkay/matcha/internal/store"
)

func Routes(store *store.Store) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/login", Authenticate(store))
	router.Post("/register", Register(store))
	router.With(middlewares.AuthTokenCtx(store), middlewares.AuthGuard(store)).Get("/logout", Logout(store))
	return router
}
