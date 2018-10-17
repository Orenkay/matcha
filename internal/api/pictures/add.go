package pictures

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api"
	"github.com/orenkay/matcha/internal/crypto"
	"github.com/orenkay/matcha/internal/store"
)

func Add(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)

		count, err := s.PicturesService.PicturesCount(user.ID)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
			if count == 5 {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("You can't upload more pictures (max 5)")))
				return
			}
		}

		file, _, err := r.FormFile("picture")
		{
			if err != nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("Invalid file")))
				return
			}
			defer file.Close()
		}

		fileBytes, err := ioutil.ReadAll(file)
		{
			if err != nil {
				render.Render(w, r, api.ErrInvalidRequest(errors.New("Invalid file")))
				return
			}
		}

		filetype := http.DetectContentType(fileBytes)
		switch filetype {
		case "image/jpeg", "image/jpg":
			break
		default:
			render.Render(w, r, api.ErrInvalidRequest(errors.New("Invalid image type, only .jpeg image is accepted")))
			return
		}

		token, err := crypto.RandToken(16)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		filename := fmt.Sprintf("%d_%s", user.ID, token)
		path := path.Join(os.Getenv("MATCHA_PATH"), "assets/uploads", filename+".jpg")
		newFile, err := os.Create(path)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}

		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}

		p := &store.Picture{
			UserID: user.ID,
			Path:   "http://" + r.Host + "/pictures/" + filename,
			IsPP:   false,
		}
		if err := s.PicturesService.Add(p); err != nil {
			render.Render(w, r, api.ErrInternal(err))
			return
		}
		render.Render(w, r, api.DefaultResponse(http.StatusCreated, p))
	}
}
