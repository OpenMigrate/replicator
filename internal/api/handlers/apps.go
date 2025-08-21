package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"replicator/internal/api/dto"
	mw "replicator/internal/api/middleware"
	"replicator/internal/models"
	"replicator/internal/storage"
)

type createAppReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type addServersReq struct {
	ServerIDs []string `json:"metadata_ids"`
}

// POST /api/apps
func CreateAppHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil {
		log.Error("CreateAppHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	var req createAppReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("CreateAppHandler: decode failed", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	var cnt int64
	if err := store.DB.Model(&models.App{}).Where("name = ?", name).Count(&cnt).Error; err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if cnt > 0 {
		http.Error(w, "name already exists", http.StatusConflict)
		return
	}

	app, err := store.CreateApp(storage.AppCreate{
		ID:          uuid.NewString(),
		Name:        name,
		Description: req.Description,
	})
	if err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	resp := dto.App{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// DELETE /api/apps/{id}
func DeleteAppHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil || store.DB == nil {
		log.Error("DeleteAppHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	if err := store.DeleteApp(storage.AppSelector{ID: &id}); err != nil {
		log.Error("DeleteAppHandler: db error", "error", err.Error())
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.Status{Status: "ok"})
}

// GET /api/apps
func ListAppsHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil {
		log.Error("ListAppsHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	afterID := r.URL.Query().Get("after_id")
	limit := 50
	if lq := r.URL.Query().Get("limit"); lq != "" {
		if v, err := strconv.Atoi(lq); err == nil && v > 0 && v <= 500 {
			limit = v
		}
	}

	items, next, err := store.ListApps(afterID, limit)
	if err != nil {
		http.Error(w, "list failed", http.StatusInternalServerError)
		return
	}

	out := dto.AppList{
		NextCursor: next,
		Items:      make([]dto.App, 0, len(items)),
	}
	for i := range items {
		out.Items = append(out.Items, dto.App{
			ID:          items[i].ID,
			Name:        items[i].Name,
			Description: items[i].Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// GET /api/apps/{id}
func GetAppByIDHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil {
		log.Error("GetAppByIDHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	app, err := store.FindApp(storage.AppSelector{ID: &id})
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := dto.App{
		ID:          app.ID,
		Name:        app.Name,
		Description: app.Description,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// POST /api/apps/{appID}/servers
func AddServersToAppHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil {
		log.Error("AddServersToAppHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	appID := chi.URLParam(r, "appID")

	var req addServersReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("AddServersToAppHandler: decode failed", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.ServerIDs) == 0 {
		http.Error(w, "metadata_ids required", http.StatusBadRequest)
		return
	}

	if err := store.ModifyAppServers(
		storage.AppSelector{ID: &appID},
		req.ServerIDs,
		storage.MembershipAdd); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	resp := dto.StatusCount{Status: "ok", Count: len(req.ServerIDs)}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// DELETE /api/apps/{appID}/servers/{serverID}
func RemoveServerFromAppHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil {
		log.Error("RemoveServerFromAppHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	appID := chi.URLParam(r, "appID")
	serverID := chi.URLParam(r, "serverID")

	if err := store.ModifyAppServers(
		storage.AppSelector{ID: &appID},
		[]string{serverID},
		storage.MembershipRemove); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.Status{Status: "ok"})
}

// GET /api/apps/{appID}/servers
func ListServersForAppHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil {
		log.Error("ListServersForAppHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	appID := chi.URLParam(r, "appID")
	afterID := r.URL.Query().Get("after_id")
	limit := 50
	if lq := r.URL.Query().Get("limit"); lq != "" {
		if v, err := strconv.Atoi(lq); err == nil && v > 0 && v <= 500 {
			limit = v
		}
	}

	servers, total, next, err := store.ListAppServers(
		storage.AppSelector{ID: &appID},
		storage.Cursor{AfterID: afterID, Limit: limit},
	)
	if err != nil {
		http.Error(w, "list failed", http.StatusInternalServerError)
		return
	}

	out := dto.ServerList{
		Total:      total,
		NextCursor: next,
		Items:      make([]dto.Server, 0, len(servers)),
	}
	for i := range servers {
		out.Items = append(out.Items, dto.Server{
			ID:           servers[i].ID,
			Hostname:     servers[i].Hostname,
			OS:           servers[i].OS,
			Arch:         servers[i].Arch,
			NumCPU:       servers[i].NumCPU,
			TimestampUTC: servers[i].TimestampUTC,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
