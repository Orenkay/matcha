package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/mail"
	"github.com/orenkay/matcha/internal/store"
)

type ResetRequest struct {
	Email string `json:"email"`
}

func (data *ResetRequest) Bind(r *http.Request) error {
	return nil
}

func ResetPassword(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &ResetRequest{}
		user, err := s.UserService.UserByEmail(data.Email)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if user == nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("There is no user associated with the provided email address")))
				return
			}
		}
		if err := mail.SendPassReset(s, user); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		if err := s.AuthTokenService.DeleteByUserID(user.ID); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, "Password reseted"))
	}
}
