package handlers

import (
	"encoding/json"
	"net/http"

	mw "replicator/internal/api/middleware"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// ListServersHandler responds with the list of all servers stored in the backend.
//
// It retrieves the storage instance and logger from the request context.
// If storage is missing or the list operation fails, it returns an HTTP 500.
// On success, it encodes the server list as JSON.
func ListServersHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)

	storage := mw.StoreFrom(r)
	if storage == nil {
		log.Error("ListServersHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	data, err := storage.ListServers()
	if err != nil {
		log.Error("ListServersHandler: ListServers failed", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error("ListServersHandler: encode failed", "error", err.Error())
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

// GetServerHandler responds with the metadata for a specific server by ID.
//
// It retrieves the storage instance and logger from the request context.
// If the server is not found, it returns HTTP 404.
// If storage is missing or a retrieval/encoding error occurs, it returns HTTP 500.
func GetServerHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)

	storage := mw.StoreFrom(r)
	if storage == nil {
		log.Error("GetServerHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	log.Debug("GetServerHandler: fetching server", "id", id)

	md, err := storage.GetServer(id)
	if err == gorm.ErrRecordNotFound {
		log.Warn("GetServerHandler: not found", "id", id)
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Error("GetServerHandler: GetServer failed", "id", id, "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if md.ID == "" {
		log.Warn("GetServerHandler: empty result", "id", id)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(md); err != nil {
		log.Error("GetServerHandler: encode failed", "id", id, "error", err.Error())
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

