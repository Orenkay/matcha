package profiles

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

type CreateRequest struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Gender     string `json:"gender"`
	Attraction string `json:"attraction"`
	Bio        string `json:"bio"`
}

func (data *CreateRequest) Bind(r *http.Request) error {
	// todo: Validation
	return nil
}

// Create handle POST /profiles requests
func Create(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		data := &CreateRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		profile, err := s.ProfileService.Profile(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if profile != nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("profile already created")))
				return
			}
		}
		profile = &store.Profile{
			UserID:     user.ID,
			LastName:   data.LastName,
			FirstName:  data.FirstName,
			Gender:     data.Gender,
			Attraction: data.Attraction,
			Bio:        data.Bio,
		}
		if err := s.ProfileService.Add(profile); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusCreated, profile))
	}
}
