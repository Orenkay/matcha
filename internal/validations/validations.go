package validations

import (
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

func Username(ve *api.ValidationError, value string) {
	if len(value) < 4 || len(value) > 32 {
		ve.Add("user", "Username length must be between 4 and 32 included.")
	}
}

func Password(ve *api.ValidationError, value string) {
	if len(value) < 4 || len(value) > 32 {
		ve.Add("pass", "Password length must be between 4 and 32 included.")
	}
}

func EmailTaken(ve *api.ValidationError, s store.UserService, value string) {
	user, err := s.UserByEmail(value)
	{
		if err != nil {
			ve.InternalError(err)
		}
		if user != nil {
			ve.Add("email", "Already taken")
		}
	}
}

func UsernameTaken(ve *api.ValidationError, s store.UserService, value string) {
	user, err := s.UserByName(value)
	{
		if err != nil {
			ve.InternalError(err)
		}
		if user != nil {
			ve.Add("user", "Already taken")
		}
	}
}
