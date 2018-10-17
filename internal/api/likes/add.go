package likes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Add(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}

		if user.ID == int64(id) {
			render.Render(w, r, api.ErrInvalidRequest(errors.New("Wow you tried to like yourself ? o.o")))
			return
		}

		like, err := s.LikesService.Like(user.ID, int64(id))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if like {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("You already like this profile")))
				return
			}
		}
		if err := s.LikesService.Add(user.ID, int64(id)); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		match, err := s.LikesService.Like(int64(id), user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if match {
				if err := s.MatchService.Add(user.ID, int64(id)); err != nil {
					render.Render(w, r, api.ErrInternal(err))
					return
				}
			}
		}
		render.Render(w, r, api.DefaultResponse(http.StatusCreated, nil))
	}
}
