package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/store"
)

type AccountEditRequest struct {
	Email           string `json:"email"`
	Username        string `json:"user"`
	CurrentPassword string `json:"currPass"`
}

func (data *AccountEditRequest) Bind(r *http.Request) error {
	// ve := &api.ValidationError{}
	// ve.Validation.Source = "AccountEdit"
	return nil
}

func EditAccount(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)

		data := &AccountEditRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrValidation(err))
			return
		}

		if err := crypto.ComparePassword(user.Password, data.CurrentPassword); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(errors.New("Current password is invalid")))
			return
		}

		user.Email = data.Email
		user.Username = data.Username
		if err := s.UserService.Update(user); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	})
}
