package messages

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
	"github.com/orenkay/matcha/internal/validations"
)

type SendRequest struct {
	Message string `json:"message"`
}

func (data *SendRequest) Bind(r *http.Request) error {
	ve := &api.ValidationError{}
	ve.Validation.Source = "messages"

	validations.Message(ve, data.Message)

	if ve.Len() > 0 {
		return ve
	}

	return nil
}

func SendMessage(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		target, err := strconv.Atoi(chi.URLParam(r, "to"))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		match, err := s.MatchService.Match(user.ID, int64(target))
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if !match {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("You need to match before sending messages :(")))
				return
			}
		}
		data := &SendRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}
		msg := &store.Message{
			Sender:   user.ID,
			Reciever: int64(target),
			Message:  data.Message,
			Date:     time.Now().Unix(),
		}
		if err := s.MessageService.Add(msg); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusCreated, msg))
	}
}
