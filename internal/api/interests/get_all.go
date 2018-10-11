package interests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func GetAll(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		interests, err := s.InterestService.AllBySlug(slug, 10)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, interests))
	}
}
