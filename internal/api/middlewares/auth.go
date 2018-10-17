package middlewares

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func AuthTokenCtx(s *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := s.AuthTokenService.Token(r.Header.Get("X-Auth-Token"))
			{
				if err != nil {
					render.Render(w, r, api.ErrInternal(err))
					return
				}
			}

			ctx := context.WithValue(r.Context(), "auth-token", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthGuard(s *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := r.Context().Value("auth-token").(*jwt.Token)
			{
				if !ok || token == nil {
					render.Render(w, r, api.ErrAuthenticate())
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
