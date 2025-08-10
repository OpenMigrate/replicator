package api

import (
	"encoding/json"
	"net/http"
	"replicator/internal/storage"

	"github.com/go-chi/chi/v5"
)

func ListServersHandler(w http.ResponseWriter, r *http.Request) {
  data := storage.ListServers()
  w.Header().Set("content-type", "application/json")
  json.NewEncoder(w).Encode(data)
}

func GetServerHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	md, ok := storage.GetServer(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(md)
}
