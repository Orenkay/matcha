package interests

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Remove(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		value := strings.ToLower(chi.URLParam(r, "slug"))

		if err := s.InterestService.Remove(user.ID, value); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		render.Render(w, r, api.DefaultResponse(http.StatusOK, nil))
	}
}
