package users

import (
	"github.com/go-chi/chi"
	"github.com/orenkay/matcha/internal/api/middlewares"
	"github.com/orenkay/matcha/internal/store"
)

func Routes(store *store.Store) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middlewares.AuthTokenCtx(store))

	router.Route("/me", func(r chi.Router) {
		r.Use(middlewares.UserMeCtx(store))
		r.Get("/", Profile(store))
		r.Patch("/password", EditPassword(store))
		r.Patch("/", EditAccount(store))
		r.Delete("/{pass}", DeleteAccount(store))
	})

	router.Route("/{userID}", func(r chi.Router) {
		r.Use(middlewares.UserCtx(store))
		r.Get("/", Profile(store))
		r.Get("/validate/{code}", Validate(store))
		r.Get("/validate/resend", ResendValidation(store))
	})

	return router
}
