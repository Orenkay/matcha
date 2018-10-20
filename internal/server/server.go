package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/orenkay/matcha/internal/api/auth"
	"github.com/orenkay/matcha/internal/api/blocks"
	"github.com/orenkay/matcha/internal/api/interests"
	"github.com/orenkay/matcha/internal/api/likes"
	"github.com/orenkay/matcha/internal/api/localisations"
	"github.com/orenkay/matcha/internal/api/matcher"
	"github.com/orenkay/matcha/internal/api/messages"
	"github.com/orenkay/matcha/internal/api/notifications"
	"github.com/orenkay/matcha/internal/api/pictures"
	"github.com/orenkay/matcha/internal/api/profiles"
	"github.com/orenkay/matcha/internal/api/reports"
	"github.com/orenkay/matcha/internal/api/users"
	"github.com/orenkay/matcha/internal/store"
)

type Server struct {
	store  *store.Store
	router *chi.Mux
}

func New(store *store.Store) *Server {
	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://192.168.1.242:8080", "http://0.0.0.0:8080" /* TODO: remove me later: */, "*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Auth-Token"},
	})

	// init middlewares
	router.Use(cors.Handler)
	router.Use(middleware.RequestID)
	// router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// init routes
	router.Mount("/auth", auth.Routes(store))
	router.Mount("/users", users.Routes(store))
	router.Mount("/profiles", profiles.Routes(store))
	router.Mount("/interests", interests.Routes(store))
	router.Mount("/pictures", pictures.Routes(store))
	router.Mount("/loc", localisations.Routes(store))
	router.Mount("/likes", likes.Routes(store))
	router.Mount("/notifs", notifications.Routes(store))
	router.Mount("/messages", messages.Routes(store))
	router.Mount("/block", blocks.Routes(store))
	router.Mount("/report", reports.Routes(store))
	router.Mount("/matcher", matcher.Routes(store))
	return &Server{
		store:  store,
		router: router,
	}
}

func (s *Server) Run(addr string) error {
	fmt.Printf("matcha api server running on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
