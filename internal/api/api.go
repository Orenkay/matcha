package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	HTTPStatusCode int         `json:"-"`
	Data           interface{} `json:"data,omitempty"`
}

func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, resp.HTTPStatusCode)
	return nil
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (resp *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, 200)
	return nil
}

func NewResponse(code int, data interface{}) render.Renderer {
	return &Response{
		HTTPStatusCode: 200,
		Data:           data,
	}
}

func CodeResponse(code int) render.Renderer {
	return &Response{
		HTTPStatusCode: code,
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
	errors     int `json:"-"`
	Validation struct {
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

func (e *ValidationError) Error() string {
	return "TODO: ValidationError.Error() implementation"
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
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error.",
	}
}

func ErrUnauthorized() render.Renderer {
	return &ErrResponse{
		Err:            errors.New(""),
		HTTPStatusCode: 401,
		StatusText:     "Unauthorized.",
	}
}

func ErrValidation(err error) render.Renderer {
	ve := err.(*ValidationError)
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid query parameters.",
		Data:           ve,
	}
}
