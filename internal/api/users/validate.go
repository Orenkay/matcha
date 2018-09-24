package users

import (
	"net/http"

	"github.com/orenkay/matcha/internal/store"
)

func Validate(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// code := chi.URLParam(r, "code")
	}
}
