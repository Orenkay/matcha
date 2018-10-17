package profiles

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Profile(s *store.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		userID, err := strconv.Atoi(chi.URLParam(r, "target"))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		profile, err := s.ProfileService.Profile(int64(userID))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if profile == nil {
				render.Render(w, r, api.ErrNotFound(errors.New(`profile doens't exists`)))
				return
			}
		}
		interests, err := s.InterestService.AllByUser(int64(userID))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if len(interests) == 0 {
				render.Render(w, r, api.ErrNotFound(errors.New(`profile doens't exists`)))
				return
			}
		}
		pictures, err := s.PicturesService.Pictures(int64(userID))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			pp, err := s.PicturesService.PP(int64(userID))
			{
				if err != nil {
					render.Render(w, r, api.ErrInternal(err))
					return
				}
				if pp == nil {
					render.Render(w, r, api.ErrNotFound(errors.New(`profile doens't exists`)))
					return
				}
			}
		}
		distance, err := s.LocalisationService.DistanceBetween(user.ID, int64(userID))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if distance == -1 {
				render.Render(w, r, api.ErrNotFound(errors.New(`profile doens't exists`)))
				return
			}
		}
		likes, err := s.LikesService.Count(int64(userID))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		visites, err := s.HistoryService.Count(int64(userID))
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
			"profile": profile,
			"loc": &store.Localisation{
				Address: fmt.Sprintf("%.2f km from you", distance/1000),
			},
			"popularity": popularity,
			"interests":  interests,
			"pictures":   pictures,
		}))
	})
}
