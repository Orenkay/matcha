package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func UserCtx(s *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user *store.User

			userID := chi.URLParam(r, "userID")
			{
				if userID != "" {
					ID, err := strconv.Atoi(userID)
					if err != nil {
						render.Render(w, r, api.ErrInternal(err))
						return
					}
					user, err = s.UserService.User(int64(ID))
					if err != nil {
						render.Render(w, r, api.ErrInternal(err))
						return
					}
				}
			}

			if user == nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("Unable to restore user with id "+userID)))
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserMeCtx(s *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var user *store.User

			token, ok := r.Context().Value("auth-token").(*jwt.Token)
			{
				if !ok || token == nil {
					render.Render(w, r, api.ErrAuthenticate())
					return
				}
			}

			claims := token.Claims.(jwt.MapClaims)
			userID := int64(claims["userId"].(float64))

			user, err := s.UserService.User(userID)
			{
				if err != nil {
					render.Render(w, r, api.ErrUnauthorized(errors.New("Your auth token is prolly corrupted, please retry or relog")))
					return
				}
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
