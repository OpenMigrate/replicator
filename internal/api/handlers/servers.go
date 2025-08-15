package handlers

import (
	"encoding/json"
	"net/http"
	mw "replicator/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func ListServersHandler(w http.ResponseWriter, r *http.Request) {
	storage := mw.StoreFrom(r)
	if storage == nil {
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}
	
	data, err := storage.ListServers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
	
}

func GetServerHandler(w http.ResponseWriter, r *http.Request) {
	
	storage := mw.StoreFrom(r)
	if storage == nil {
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}
	
	id := chi.URLParam(r, "id")
	
	md, err := storage.GetServer(id)
	if err == gorm.ErrRecordNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if md.ID == "" {
		http.NotFound(w, r)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(md); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

