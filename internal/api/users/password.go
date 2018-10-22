package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/store"
	"github.com/orenkay/matcha/internal/validations"
)

type PasswordEditRequest struct {
	CurrentPassword string `json:"currPass"`
	Password        string `json:"pass"`
}

func (data *PasswordEditRequest) Bind(r *http.Request) error {
	ve := &api.ValidationError{}
	ve.Validation.Source = "register"

	validations.Password(ve, data.Password)

	if ve.Len() > 0 {
		return ve
	}

	pass, err := crypto.EncryptPassword(data.Password)
	if err != nil {
		return err
	}
	data.Password = string(pass)
	return nil
}

func EditPassword(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)

		data := &PasswordEditRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrValidation(err))
			return
		}

		if err := crypto.ComparePassword(user.Password, data.CurrentPassword); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(errors.New("Current password is invalid")))
			return
		}

		user.Password = data.Password
		if err := s.UserService.Update(user); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		s.AuthTokenService.Delete(r.Header.Get("X-Auth-Token"))
		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	})
}
