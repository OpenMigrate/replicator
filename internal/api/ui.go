package api

import (
	"html/template"
	"net/http"
	"replicator/internal/storage"

	"github.com/go-chi/chi/v5"
)

var templates = template.Must(template.ParseGlob("pkg/ui/templates/*.html"))

func IndexPage(w http.ResponseWriter, r *http.Request) {
	data := storage.ListServers()
	templates.ExecuteTemplate(w, "index.html", data)
}

func ServerPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	md, ok := storage.GetServer(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	templates.ExecuteTemplate(w, "server.html", md)
}

