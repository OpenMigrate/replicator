package ui

import (
	"html/template"
	"net/http"
	mw "replicator/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

var templates = template.Must(template.ParseGlob("pkg/ui/templates/*.html"))

func IndexPage(w http.ResponseWriter, r *http.Request) {
	storage := mw.StoreFrom(r)
	if storage == nil {
		http.Error(w, "store missing", 500)
		return
	}

	data, _ := storage.ListServers()
	templates.ExecuteTemplate(w, "index.html", data)
}

func ServerPage(w http.ResponseWriter, r *http.Request) {
	storage := mw.StoreFrom(r)
	if storage == nil {
		http.Error(w, "store missing", 500)
		return
	}

	id := chi.URLParam(r, "id")
	md, err := storage.GetServer(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	templates.ExecuteTemplate(w, "server.html", md)
}
