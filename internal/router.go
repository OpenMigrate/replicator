package internal

import (
	"net/http"
	"replicator/internal/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func NewRouter() http.Handler {
  r := chi.NewRouter();
  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)

  r.Post("/discover", api.DiscoverHandler)

  r.Route("/api", func(r chi.Router) {
    r.Post("/discover", api.DiscoverHandler)
    r.Get("/servers", api.ListServersHandler)
    r.Get("/servers/{id}", api.GetServerHandler)
  })

  // UI routes
  r.Get("/", api.IndexPage)
	r.Get("/server/{id}", api.ServerPage)

  return r
}
