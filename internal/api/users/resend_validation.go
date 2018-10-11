package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/mail"
	"github.com/orenkay/matcha/internal/store"
)

func ResendValidation(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)

		v, err := s.ValidationService.ValidationCode(user.ID)
		if err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		if v.Used {
			render.Render(w, r, api.ErrInvalidRequest(errors.New("Validation link already used")))
			return
		}

		if err := s.ValidationService.Remove(user.ID); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		err = mail.SendValidationMail(r, s, user)
		if err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	}
}
