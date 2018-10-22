package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func (resp *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, 200)
	return nil
}

type Response struct {
	HTTPStatusCode int         `json:"-"`
	Data           interface{} `json:"data,omitempty"`
}

func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, resp.HTTPStatusCode)
	return nil
}

func DefaultResponse(code int, data interface{}) render.Renderer {
	return &Response{
		HTTPStatusCode: code,
		Data:           data,
	}
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string      `json:"status"`
	AppCode    int64       `json:"code,omitempty"`
	ErrorText  string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

type ValidationErrorDetail struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type ValidationError struct {
	errors        int
	internalError *error
	Validation    struct {
		Details []ValidationErrorDetail `json:"details"`
		Source  string                  `json:"source"`
		Keys    []string                `json:"keys"`
	} `json:"validation"`
}

func (e *ValidationError) Add(path string, message string) {
	e.Validation.Keys = append(e.Validation.Keys, path)
	e.Validation.Details = append(e.Validation.Details, ValidationErrorDetail{
		Path:    path,
		Message: message,
	})
	e.errors++
}

func (e *ValidationError) Len() int {
	return e.errors
}

func (e *ValidationError) InternalError(err error) {
	e.errors++
	e.internalError = &err
}

func (e *ValidationError) Error() string {
	return "ValidationError"
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrInternal(err error) render.Renderer {
	fmt.Println(err.Error())
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error.",
	}
}

func ErrNotFound(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     err.Error(),
	}
}

func ErrAuthenticate() render.Renderer {
	return &ErrResponse{
		Err:            errors.New(""),
		HTTPStatusCode: 401,
		StatusText:     "Unauthorized place. You need to provide an auth token.",
	}
}

func ErrUnauthorized(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 403,
		StatusText:     err.Error(),
	}
}

func ErrValidation(err error) render.Renderer {
	if ve, ok := err.(*ValidationError); ok {
		if ve.internalError != nil {
			return ErrInternal(*ve.internalError)
		}
		return &ErrResponse{
			Err:            err,
			HTTPStatusCode: 400,
			StatusText:     "Invalid query parameters.",
			Data:           ve,
		}
	}
	return ErrInternal(err)
}
