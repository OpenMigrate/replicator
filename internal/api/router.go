package api

import (
	"net/http"
	"replicator/internal/storage"

	"replicator/internal/api/handlers"
	mw "replicator/internal/api/middleware"
	"replicator/internal/api/ui"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(store *storage.Store) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.WithStore(store))

	r.Post("/discover", handlers.DiscoverHandler)

	r.Route("/api", func(r chi.Router) {
		r.Post("/discover", handlers.DiscoverHandler)
		r.Get("/servers", handlers.ListServersHandler)
		r.Get("/servers/{id}", handlers.GetServerHandler)
	})

	// UI routes
	r.Get("/", ui.IndexPage)
	r.Get("/server/{id}", ui.ServerPage)

	return r
}
