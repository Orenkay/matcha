package profiles

import (
	"net/http"
	"strconv"

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
		r.Post("/", Create(s))
		r.Get("/", Profile(s))
		r.Patch("/edit", Edit(s))
	})

	router.Route("/{target}", func(r chi.Router) {
		r.Use(middlewares.UserMeCtx(s))
		r.Get("/", Profile(s))
		r.Get("/presence", Presence(s))
		r.Get("/visit", Visit(s))
	})

	return router
}

func Presence(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target, err := strconv.Atoi(chi.URLParam(r, "target"))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		last, err := s.PresenceService.LastHeartbeat(int64(target))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, last))
	}
}

func Visit(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		target, err := strconv.Atoi(chi.URLParam(r, "target"))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		if user.ID != int64(target) {
			i := &store.HistoryItem{
				UserID: int64(target),
				FromID: user.ID,
				Event:  "visite",
			}
			visited, err := s.HistoryService.Exists(i)
			{
				if err != nil {
					render.Render(w, r, api.ErrInternal(err))
					return
				}
			}
			if !visited {
				if err := s.HistoryService.Add(i); err != nil {
					render.Render(w, r, api.ErrInternal(err))
					return
				}
				s.NotificationService.Push(i.UserID, user.ID, "visite", i)
			}
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	}
}
