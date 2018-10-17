package likes

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
		r.Post("/{id}", Add(s))
		r.Delete("/{id}", Remove(s))
		r.Get("/{id}", Like(s))
	})

	return router
}

func Like(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		target, err := strconv.Atoi(chi.URLParam(r, "id"))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		like, err := s.LikesService.Like(user.ID, int64(target))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		likeMe, err := s.LikesService.Like(int64(target), user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, render.M{
			"liked":  like,
			"likeMe": likeMe,
		}))
	}
}
