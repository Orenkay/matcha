package pictures

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

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
				render.Render(w, r, api.ErrInvalidRequest(errors.New("This picture doesn't exist")))
				return
			}
		}
		if err := s.PicturesService.Delete(user.ID, p.ID); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		// here we remove locally the img
		{
			ss := strings.Split(p.Path, "/")
			filename := ss[len(ss)-1]
			os.Remove(path.Join(os.Getenv("MATCHA_PATH"), "assets/uploads", filename+".jpg"))
		}

		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	}
}
