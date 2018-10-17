package profiles

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

type EditRequest struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Birthdate  string `json:"birthDate"`
	Gender     string `json:"gender"`
	Attraction string `json:"attraction"`
	Bio        string `json:"bio"`
}

func (data *EditRequest) Bind(r *http.Request) error {
	// todo: Validation
	return nil
}

// Create handle PUT /profiles/edit requests
func Edit(s *store.Store) http.HandlerFunc {
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
		}
		profile.LastName = data.LastName
		profile.FirstName = data.FirstName
		profile.Birthdate = data.Birthdate
		profile.Gender = data.Gender
		profile.Attraction = data.Attraction
		profile.Bio = data.Bio
		if err := s.ProfileService.Update(profile); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusOK, profile))
	}
}
