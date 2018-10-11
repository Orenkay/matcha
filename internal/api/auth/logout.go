package auth

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

// Logout handle GET /auth/logout requests
func Logout(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.AuthTokenService.Delete(r.Header.Get("X-Auth-Token"))
		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	}
}
