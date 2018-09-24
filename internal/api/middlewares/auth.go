package middlewares

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/store"
)

func AuthTokenCtx(s *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				token       *jwt.Token
				tokenString = r.Header.Get("X-Auth-Token")
			)

			if s.AuthTokenService.Exists(tokenString) {
				token, err := crypto.DecodeJWT(tokenString)

				// Here we'll check if token has expired
				// If token has expired we remove it from the store
				if !token.Valid {
					if ve, ok := err.(jwt.ValidationError); ok {
						if ve.Errors&jwt.ValidationErrorExpired != 0 {
							s.AuthTokenService.Delete(tokenString)
						}
					}
				}
			}

			ctx := context.WithValue(r.Context(), "auth-token", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
