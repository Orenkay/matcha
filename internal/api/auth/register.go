package auth

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/store"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"user"`
	Password string `json:"pass"`
}

func (data *RegisterRequest) Bind(r *http.Request) error {
	ve := &api.ValidationError{}
	ve.Validation.Source = "register"

	if len(data.Username) < 4 || len(data.Username) > 32 {
		ve.Add("username", "Username length must be between 4 and 32 included.")
	}

	if len(data.Password) < 4 || len(data.Password) > 32 {
		ve.Add("password", "Password length must be between 4 and 32 included.")
	}

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

func canRegisterCheck(s *store.Store, data *RegisterRequest) error {
	ve := &api.ValidationError{}

	// We check if email is already taken
	user, err := s.UserService.UserByEmail(data.Email)
	{
		if err != nil {
			return err
		}
		if user != nil {
			ve.Add("email", "Already taken")
		}
	}

	// We check if username is already taken
	user, err = s.UserService.UserByName(data.Username)
	{
		if err != nil {
			return err
		}
		if user != nil {
			ve.Add("username", "Already taken")
		}
	}

	if ve.Len() > 0 {
		return ve
	}

	return nil
}

// Register handle POST /auth/register requests
func Register(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &RegisterRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrValidation(err))
			return
		}

		if err := canRegisterCheck(s, data); err != nil {
			if err, ok := err.(*api.ValidationError); ok {
				render.Render(w, r, api.ErrValidation(err))
				return
			}
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		err := s.UserService.Add(&store.User{
			Email:    data.Email,
			Username: data.Username,
			Password: data.Password,
		})
		if err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.CodeResponse(http.StatusCreated))
	}
}
