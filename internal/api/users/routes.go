package users

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/orenkay/matcha/internal/store"
)

func UserCtx(s *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}
}

func Routes(store *store.Store) *chi.Mux {
	router := chi.NewRouter()
	router.Get("{userID}/validate/{code}", Validate(store))
	return router
}
