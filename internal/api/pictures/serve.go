package pictures

import (
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
