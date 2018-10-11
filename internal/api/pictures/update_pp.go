package pictures

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func UpdatePP(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		p, err := s.PicturesService.Picture(user.ID, int64(id))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if p == nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("image doesn't exist")))
				return
			}
		}

		pp, err := s.PicturesService.PP(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if pp != nil {
				pp.IsPP = false
				if err := s.PicturesService.Update(pp); err != nil {
					render.Render(w, r, api.ErrInternal(err))
					return
				}
			}
		}

		p.IsPP = true
		if err := s.PicturesService.Update(p); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	}
}
