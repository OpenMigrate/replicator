package handlers

import (
	"encoding/json"
	"net/http"
	mw "replicator/internal/api/middleware"
	"replicator/internal/models"

	"github.com/google/uuid"
)

func DiscoverHandler(w http.ResponseWriter, r *http.Request) {
	var md models.Metadata
	err := json.NewDecoder(r.Body).Decode(&md)
	s := mw.StoreFrom(r)
	if s != nil {
		http.Error(w, "store missing", 500)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	md.ID = uuid.New().String()
	if err := s.SaveServer(md); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": md.ID})
}
