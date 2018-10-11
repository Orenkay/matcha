package interests

import (
	"github.com/go-chi/chi"
	"github.com/orenkay/matcha/internal/api/middlewares"
	"github.com/orenkay/matcha/internal/store"
)

func Routes(s *store.Store) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middlewares.AuthTokenCtx(s))
	router.Route("/me", func(r chi.Router) {
		r.Use(middlewares.UserMeCtx(s))
		r.Post("/", Add(s))
		r.Delete("/{slug}", Remove(s))
	})
	router.Get("/", GetAll(s))
	router.Get("/{slug}", GetAll(s))
	return router
}
