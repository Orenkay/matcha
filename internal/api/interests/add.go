package interests

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/store"
)

type AddRequest struct {
	Value string `json:"value"`
}

func (data *AddRequest) Bind(r *http.Request) error {
	data.Value = strings.ToLower(data.Value)
	return nil
}

func Add(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		data := &AddRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		i, err := s.InterestService.Interest(user.ID, data.Value)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if i != nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("Duplicate interests")))
				return
			}
		}
		i = &store.Interest{
			UserID: user.ID,
			Value:  data.Value,
		}
		err = s.InterestService.Add(i)
		if err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusCreated, i))
	}
}
