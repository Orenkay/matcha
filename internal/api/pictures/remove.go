package pictures

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"

	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Remove(s *store.Store) http.HandlerFunc {
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
				render.Render(w, r, api.ErrInvalidRequest(errors.New("The pictures doesn't exist")))
				return
			}
		}
		if err := s.PicturesService.Delete(user.ID, p.ID); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		// here we dont handle error: see https://github.com/Orenkay/matcha/blob/master/internal/api/pictures/serve.go#L25
		os.Remove(path.Join(os.Getenv("MATCHA_PATH"), "assets/uploads", p.Path+".jpg"))
		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	}
}
