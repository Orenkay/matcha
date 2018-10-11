package pictures

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
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

		// check file type, detectcontenttype only needs the first 512 bytes
		filetype := http.DetectContentType(fileBytes)
		switch filetype {
		case "image/jpeg", "image/jpg", "image/png":
			break
		default:
			render.Render(w, r, api.ErrInvalidRequest(errors.New("Invalid file type")))
			return
		}

		fileExt, err := mime.ExtensionsByType(filetype)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}

		token, err := crypto.RandToken(16)
		{
			if err != nil {
				render.Render(w, r, api.ErrInternal(err))
				return
			}
		}
		filename := path.Join(os.Getenv("MATCHA_PATH"), "assets/uploads", fmt.Sprintf("%d_%s%s", user.ID, token, fileExt[0]))
		newFile, err := os.Create(filename)
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

		render.Render(w, r, api.DefaultResponse(http.StatusCreated, nil))
	}
}
