package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"

	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Validate(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		validated, err := s.ValidationService.IsValidated(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if validated {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("account is already activated")))
				return
			}
		}

		code := chi.URLParam(r, "code")
		valid, err := s.ValidationService.CheckCode(user.ID, code)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}

			if !valid {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("provided code is invalid")))
				return
			}
		}
		if err := s.ValidationService.Validate(user.ID, code); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		render.Render(w, r, api.DefaultResponse(http.StatusOK, "account is now activated"))
	}
}
