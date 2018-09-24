package users

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
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

func Routes(store *store.Store) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/{userID}", func(r chi.Router) {
		r.Use(UserCtx(store))
		r.Get("/validate/{code}", Validate(store))
	})

	return router
}
