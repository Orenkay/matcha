package users

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Profile(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		profile, err := s.ProfileService.Profile(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		loc, err := s.LocalisationService.Localisation(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		interests, err := s.InterestService.AllByUser(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		pictures, err := s.PicturesService.Pictures(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, render.M{
			"account":   user,
			"profile":   profile,
			"loc":       loc,
			"interests": interests,
			"pictures":  pictures,
		}))
	})
}
