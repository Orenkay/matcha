package localisations

import (
	"errors"
	"net/http"

	"github.com/orenkay/matcha/internal/localisation"

	"github.com/go-chi/render"

	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

type AddRequest struct {
	PlaceID string `json:"placeId"`
}

func (data *AddRequest) Bind(r *http.Request) error {
	return nil
}

func Add(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)

		data := &AddRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		if data.PlaceID == "" {
			render.Render(w, r, api.ErrInvalidRequest(errors.New("Invalid parameters")))
			return
		}

		l, err := s.LocalisationService.Localisation(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if l != nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("You have already an location")))
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

		l = &store.Localisation{
			UserID:  user.ID,
			Lat:     p.Lat,
			Lng:     p.Lng,
			Address: p.Address,
		}
		if err := s.LocalisationService.Add(l); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, l))
	}
}
