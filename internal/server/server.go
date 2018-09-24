package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api/auth"
	"github.com/orenkay/matcha/internal/api/users"
	"github.com/orenkay/matcha/internal/store"
)

type Server struct {
	store  *store.Store
	router *chi.Mux
}

func New(store *store.Store) *Server {
	router := chi.NewRouter()

	// init middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// init routes
	router.Mount("/auth", auth.Routes(store))
	router.Mount("/users", users.Routes(store))

	return &Server{
		store:  store,
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	fmt.Printf("matcha api server running on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
