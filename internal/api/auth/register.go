package auth

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/mail"
	"github.com/orenkay/matcha/internal/store"
	"github.com/orenkay/matcha/internal/validations"
)

type RegisterRequest struct {
	s store.UserService

	Email    string `json:"email"`
	Username string `json:"user"`
	Password string `json:"pass"`
}

func (data *RegisterRequest) Bind(r *http.Request) error {
	ve := &api.ValidationError{}
	ve.Validation.Source = "register"

	data.Username = strings.TrimSpace(data.Username)
	data.Password = strings.TrimSpace(data.Password)

	validations.Username(ve, data.Username)
	validations.Password(ve, data.Password)
	validations.EmailFormat(ve, data.Email)
	validations.EmailTaken(ve, data.s, data.Email)
	validations.UsernameTaken(ve, data.s, data.Username)

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

// Register handle POST /auth/register requests
func Register(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &RegisterRequest{s: s.UserService}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrValidation(err))
			return
		}
		user := &store.User{
			Email:    data.Email,
			Username: data.Username,
			Password: data.Password,
		}
		err := s.UserService.Add(user)
		if err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		err = mail.SendValidationMail(r, s, user)
		if err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		if err := s.PresenceService.Add(user.ID); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusCreated, nil))
	}
}
