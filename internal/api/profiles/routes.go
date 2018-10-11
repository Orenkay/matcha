package profiles

import (
	"github.com/go-chi/chi"
	"github.com/orenkay/matcha/internal/api/middlewares"
	"github.com/orenkay/matcha/internal/store"
)

func Routes(s *store.Store) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middlewares.AuthTokenCtx(s))

	router.Route("/", func(r chi.Router) {
		r.Use(middlewares.UserMeCtx(s))
		r.Post("/", Create(s))
		r.Put("/edit", Edit(s))
	})
	return router
}
