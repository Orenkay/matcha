package users

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Account(s *store.Store) http.HandlerFunc {
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
		matches, err := s.MatchService.Matches(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		likes, err := s.LikesService.Count(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		visites, err := s.HistoryService.Count(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		var popularity float64
		if visites > 0 {
			popularity = float64(likes) / float64(visites)
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, render.M{
			"account": user,
			"meta": render.M{
				"profile":    profile,
				"loc":        loc,
				"interests":  interests,
				"pictures":   pictures,
				"matches":    matches,
				"popularity": popularity,
			},
		}))
	})
}
