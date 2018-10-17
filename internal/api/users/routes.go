package users

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/api/middlewares"
	"github.com/orenkay/matcha/internal/store"
)

func Routes(s *store.Store) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middlewares.AuthTokenCtx(s))

	router.Route("/me", func(r chi.Router) {
		r.Use(middlewares.UserMeCtx(s))
		r.Get("/", Account(s))
		r.Get("/history", History(s))
		r.Patch("/password", EditPassword(s))
		r.Patch("/", EditAccount(s))
		r.Delete("/{pass}", DeleteAccount(s))
		r.Patch("/heartbeat", Heartbeat(s))
	})

	router.Route("/{userID}", func(r chi.Router) {
		r.Use(middlewares.UserCtx(s))
		r.Get("/validate/{code}", Validate(s))
		r.Get("/validate/resend", ResendValidation(s))
	})

	return router
}

func History(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		history, err := s.HistoryService.History(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, history))
	}
}
func Heartbeat(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		if err := s.PresenceService.Heartbeat(user.ID); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	}
}
