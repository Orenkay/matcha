package pictures

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Serve(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image := chi.URLParam(r, "slug")
		filename := path.Join(os.Getenv("MATCHA_PATH"), "assets/uploads", image+".jpg")
		img, err := os.Open(filename)
		{
			if os.IsNotExist(err) {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("requested image doesnt exist")))
				// Here we try to cleanup database from corrupted file
				s.PicturesService.DeleteByPath(image)
				return
			}
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		defer img.Close()
		w.Header().Set("Content-Type", "image/jpeg")
		io.Copy(w, img)
	}
}
