package localisations

import (
	"errors"
	"net/http"

	"github.com/orenkay/matcha/internal/localisation"

	"github.com/go-chi/render"

	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Update(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)

		data := &AddRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		l, err := s.LocalisationService.Localisation(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if l == nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("You don't have location xoxo")))
				return
			}
		}
		p, err := localisation.PlaceByID(data.PlaceID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		l.Lat = p.Lat
		l.Lng = p.Lng
		l.Address = p.Address
		if err := s.LocalisationService.Update(l); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, l))
	}
}
