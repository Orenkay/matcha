package mail

import (
	"net/http"
	"strconv"

	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/store"
	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

var mg = mailgun.NewMailgun("spopieul.me", "de9dfa1d53a4712cbba29bbda3d3dc18-b0aac6d0-adbbe185", "pubkey-f4fc29fbf1a8ed6ea76b074b87e5288d")

func SendMail(subject, text string, to ...string) error {
	m := mg.NewMessage(
		"Matcha <mailgun@spopieul.me>",
		subject,
		text,
		to...,
	)
	_, _, err := mg.Send(m)
	return err
}

func SendValidationMail(r *http.Request, s *store.Store, user *store.User) error {
	code, err := crypto.RandToken(32)
	{
		if err != nil {
			return err
		}
	}

	if err := s.ValidationService.Add(user.ID, code); err != nil {
		return err
	}

	link := "http://" + r.Host + "/users/" + strconv.Itoa(int(user.ID)) + "/validate/" + code
	return SendMail("Validation", link, user.Email)
}
