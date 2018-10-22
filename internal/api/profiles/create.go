package profiles

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
	"github.com/orenkay/matcha/internal/validations"
)

type CreateRequest struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	Birthdate  int64  `json:"birthDate"`
	Gender     string `json:"gender"`
	Attraction string `json:"attraction"`
	Bio        string `json:"bio"`
}

func (data *CreateRequest) Bind(r *http.Request) error {
	ve := &api.ValidationError{}
	ve.Validation.Source = "edit_profile"

	data.LastName = strings.TrimSpace(data.LastName)
	data.FirstName = strings.TrimSpace(data.FirstName)
	data.Bio = strings.TrimSpace(data.Bio)

	validations.EmptyStr(ve, "lastName", data.LastName)
	validations.EmptyStr(ve, "firstName", data.FirstName)
	validations.EmptyStr(ve, "bio", data.Bio)
	validations.Max(ve, "lastName", data.LastName, 32)
	validations.Max(ve, "firstName", data.FirstName, 32)
	validations.Birthdate(ve, data.Birthdate)
	validations.Gender(ve, data.Gender)
	validations.Attraction(ve, data.Attraction)
	validations.Max(ve, "bio", data.Bio, 200)

	if ve.Len() > 0 {
		return ve
	}

	return nil
}

// Create handle POST /profiles requests
func Create(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		data := &CreateRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrValidation(err))
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
			Birthdate:  data.Birthdate,
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
