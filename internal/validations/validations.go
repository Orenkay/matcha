package validations

import (
	"fmt"
	"regexp"
	"time"

	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

var (
	emailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func Username(ve *api.ValidationError, value string) {
	if len(value) < 4 || len(value) > 32 {
		ve.Add("user", "Username length must be between 4 and 32 included.")
	}
}

func Message(ve *api.ValidationError, value string) {
	if len(value) > 200 {
		ve.Add("message", "message length must be less than 200 characters long.")
	}
}

func Password(ve *api.ValidationError, value string) {
	if len(value) < 4 || len(value) > 32 {
		ve.Add("pass", "Password length must be between 4 and 32 included.")
	}
}

func EmptyStr(ve *api.ValidationError, field string, value string) {
	if len(value) == 0 {
		ve.Add(field, fmt.Sprintf("field is required"))
	}
}

func Max(ve *api.ValidationError, field string, value string, max int) {
	if len(value) > max {
		ve.Add(field, fmt.Sprintf("length must be %d characters long max.", max))
	}
}

func Attraction(ve *api.ValidationError, value string) {
	if value != "bi" && value != "hetero" && value != "homo" {
		ve.Add("attraction", "Invalid attraction")
	}
}

func Gender(ve *api.ValidationError, value string) {
	if value != "male" && value != "female" {
		ve.Add("gender", "Invalid gender")
	}
}

func Birthdate(ve *api.ValidationError, value int64) {
	t := time.Now()
	max := time.Date(t.Year()-18, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()).Unix()
	min := ^(1 << 31) + 1
	if value < int64(min) || value > int64(max) {
		ve.Add("birthdate", "Invalid")
	}
}

func EmailFormat(ve *api.ValidationError, value string) {
	if !emailRegexp.MatchString(value) {
		ve.Add("email", "Email is not valid.")
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
