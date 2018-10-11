package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/store"
)

func DeleteAccount(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		currPassword := chi.URLParam(r, "pass")
		if err := crypto.ComparePassword(user.Password, currPassword); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(errors.New("Current password is invalid")))
			return
		}

		if err := s.UserService.Remove(user); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		s.AuthTokenService.Delete(r.Header.Get("X-Auth-Token"))
		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	})
}
