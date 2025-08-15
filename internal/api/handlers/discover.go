package handlers

import (
	"encoding/json"
	"net/http"
	mw "replicator/internal/api/middleware"
	"replicator/internal/models"

	"github.com/google/uuid"
)

// DiscoverHandler responds with acceptiong the metadata from agent and returns it's store id
//
// It retrieves the storage instance and logger from the request context.
func DiscoverHandler(w http.ResponseWriter, r *http.Request) {
	var md models.Metadata
	log := mw.GetLogFromCtx(r)

	if err := json.NewDecoder(r.Body).Decode(&md); err != nil {
		log.Error("DiscoverHandler: JSON encoding error", "msg", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := mw.StoreFrom(r)
	if s == nil {
		log.Error("ListServersHandler: store is nil")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	md.ID = uuid.New().String()
	if err := s.SaveServer(md); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"id": md.ID}); err != nil {
		log.Error("DiscoverHandler: encode failed", "error", err.Error())
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}

}
