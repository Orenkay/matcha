package auth

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/localisation"
	"github.com/orenkay/matcha/internal/store"
	"github.com/orenkay/matcha/internal/validations"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"user"`
	Password string `json:"pass"`
	Loc      *localisation.Place
}

func (data *RegisterRequest) Bind(r *http.Request) error {
	ve := &api.ValidationError{}
	ve.Validation.Source = "register"

	validations.Username(ve, data.Username)
	validations.Password(ve, data.Password)

	if ve.Len() > 0 {
		return ve
	}

	pass, err := crypto.EncryptPassword(data.Password)
	if err != nil {
		return err
	}
	data.Password = string(pass)

	// Locate user by his IP: IDK if i should keep it since now we ask for user location
	// data.Loc, err = localisation.PlaceByIP(r.RemoteAddr)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Register handle POST /auth/register requests
func Register(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &RegisterRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrValidation(err))
			return
		}

		ve := &api.ValidationError{}
		{
			validations.EmailTaken(ve, s.UserService, data.Email)
			validations.UsernameTaken(ve, s.UserService, data.Username)
			if ve.Len() > 0 {
				render.Render(w, r, api.ErrValidation(ve))
				return
			}
		}
		user := &store.User{
			Email:    data.Email,
			Username: data.Username,
			Password: data.Password,
		}
		err := s.UserService.Add(user)
		if err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		// err = s.LocalisationService.Add(&store.Localisation{
		// 	UserID:  user.ID,
		// 	Lat:     data.Loc.Lat,
		// 	Lng:     data.Loc.Lng,
		// 	Address: data.Loc.Address,
		// })
		// if err != nil {
		// 	render.Render(w, r, api.ErrInternal(err))
		// 	return
		// }

		// err = mail.SendValidationMail(r, s, user)
		// if err != nil {
		// 	render.Render(w, r, api.ErrInternal(err))
		// 	return
		// }
		render.Render(w, r, api.DefaultResponse(http.StatusCreated, nil))
	}
}
