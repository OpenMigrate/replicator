package api

import (
	"encoding/json"
	"net/http"
	"replicator/internal/models"
	"replicator/internal/storage"

	"github.com/google/uuid"
)

func DiscoverHandler(w http.ResponseWriter, r *http.Request){
  var md models.Metadata
  err := json.NewDecoder(r.Body).Decode(&md);

  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  md.ID = uuid.New().String()
  storage.SaveServer(md)

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]string{"id":md.ID})
}
