package handlers

import (
	"encoding/json"
	"net/http"

	mw "replicator/internal/api/middleware"
)

func SeedHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil {
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	if err := store.SeedSampleData(r.Context()); err != nil {
		log.Error("SeedHandler failed", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"msg":    "Sample data inserted",
	})
}
