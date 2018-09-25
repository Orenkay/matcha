package auth

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/store"
)

type LoginRequest struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

func (data *LoginRequest) Bind(r *http.Request) error {
	return nil
}

// Authenticate handle POST /auth/login requests
func Authenticate(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &LoginRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		user, err := s.UserService.UserByName(data.Username)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if user == nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("Username / Password invalid")))
				return
			}
			if err := crypto.ComparePassword(user.Password, data.Password); err != nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("Username / Password invalid")))
				return
			}
		}

		claims := jwt.MapClaims{"userId": user.ID}
		ss, token, err := crypto.CreateJWT(claims)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		s.AuthTokenService.Add(ss, token)
		render.Render(w, r, &api.TokenResponse{Token: ss})
	}
}
