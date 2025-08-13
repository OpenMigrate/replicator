package handlers

import (
	"encoding/json"
	"net/http"
	mw "replicator/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

func ListServersHandler(w http.ResponseWriter, r *http.Request) {
	storage := mw.StoreFrom(r)
	if storage != nil {
		http.Error(w, "store missing", 500)
		return
	}

	data, err := storage.ListServers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func GetServerHandler(w http.ResponseWriter, r *http.Request) {
	storage := mw.StoreFrom(r)
	if storage != nil {
		http.Error(w, "store missing", 500)
		return
	}

	id := chi.URLParam(r, "id")
	md, err := storage.GetServer(id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if md.ID == "" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(md)
}
