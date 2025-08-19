// handlers/apps.go
package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	mw "replicator/internal/api/middleware"
	"replicator/internal/models"
	"replicator/internal/storage"
)

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugRe.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

type createAppReq struct {
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
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
	identifier := slugify(req.Identifier)
	if identifier == "" {
		http.Error(w, "identifier is required", http.StatusBadRequest)
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
	if err := store.DB.Model(&models.App{}).Where("identifier = ?", identifier).Count(&cnt).Error; err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if cnt > 0 {
		http.Error(w, "identifier already exists", http.StatusConflict)
		return
	}

	app, err := store.CreateApp(r.Context(), storage.AppCreate{
		ID:          uuid.NewString(),
		Name:        name,
		Identifier:  identifier,
		Description: req.Description,
	})
	if err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"id":          app.ID,
		"name":        app.Name,
		"identifier":  app.Identifier,
		"description": app.Description,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
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

	items, next, err := store.ListApps(r.Context(), afterID, limit)
	if err != nil {
		http.Error(w, "list failed", http.StatusInternalServerError)
		return
	}

	type appItem struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Identifier  string `json:"identifier"`
		Description string `json:"description"`
	}
	out := struct {
		NextCursor string    `json:"next_cursor"`
		Items      []appItem `json:"items"`
	}{NextCursor: next, Items: make([]appItem, 0, len(items))}

	for i := range items {
		out.Items = append(out.Items, appItem{
			ID:          items[i].ID,
			Name:        items[i].Name,
			Identifier:  items[i].Identifier,
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
	app, err := store.FindApp(r.Context(), storage.AppSelector{ID: &id})
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := map[string]string{
		"id":          app.ID,
		"name":        app.Name,
		"identifier":  app.Identifier,
		"description": app.Description,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// GET /api/apps/by-identifier/{identifier}
func GetAppByIdentifierHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil {
		log.Error("GetAppByIdentifierHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	raw := chi.URLParam(r, "identifier")
	idf := slugify(raw)

	app, err := store.FindApp(r.Context(), storage.AppSelector{Identifier: &idf})
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := map[string]string{
		"id":          app.ID,
		"name":        app.Name,
		"identifier":  app.Identifier,
		"description": app.Description,
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

	if err := store.ModifyAppServers(r.Context(),
		storage.AppSelector{ID: &appID},
		req.ServerIDs,
		storage.MembershipAdd); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"status": "ok",
		"count":  len(req.ServerIDs),
	}
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

	if err := store.ModifyAppServers(r.Context(),
		storage.AppSelector{ID: &appID},
		[]string{serverID},
		storage.MembershipRemove); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
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
		r.Context(),
		storage.AppSelector{ID: &appID},
		storage.Cursor{AfterID: afterID, Limit: limit},
	)
	if err != nil {
		http.Error(w, "list failed", http.StatusInternalServerError)
		return
	}

	type serverItem struct {
		ID           string `json:"id"`
		Hostname     string `json:"hostname"`
		OS           string `json:"os"`
		Arch         string `json:"arch"`
		NumCPU       int    `json:"num_cpu"`
		TimestampUTC string `json:"timestamp_utc"`
	}

	out := struct {
		Total      int64        `json:"total"`
		NextCursor string       `json:"next_cursor"`
		Items      []serverItem `json:"items"`
	}{Total: total, NextCursor: next, Items: make([]serverItem, 0, len(servers))}

	for i := range servers {
		out.Items = append(out.Items, serverItem{
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

// GET /api/servers?app=<identifier>  (alias)
func ListServersByAppIdentifierAliasHandler(w http.ResponseWriter, r *http.Request) {
	log := mw.GetLogFromCtx(r)
	store := mw.StoreFrom(r)
	if store == nil {
		log.Error("ListServersByAppIdentifierAliasHandler: store missing")
		http.Error(w, "store missing", http.StatusInternalServerError)
		return
	}

	idf := slugify(r.URL.Query().Get("app"))
	if idf == "" {
		http.Error(w, "app is required", http.StatusBadRequest)
		return
	}

	afterID := r.URL.Query().Get("after_id")
	limit := 50
	if lq := r.URL.Query().Get("limit"); lq != "" {
		if v, err := strconv.Atoi(lq); err == nil && v > 0 && v <= 500 {
			limit = v
		}
	}

	app, err := store.FindApp(r.Context(), storage.AppSelector{Identifier: &idf})
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	servers, total, next, err := store.ListAppServers(
		r.Context(),
		storage.AppSelector{ID: &app.ID},
		storage.Cursor{AfterID: afterID, Limit: limit},
	)
	if err != nil {
		http.Error(w, "list failed", http.StatusInternalServerError)
		return
	}

	type serverItem struct {
		ID           string `json:"id"`
		Hostname     string `json:"hostname"`
		OS           string `json:"os"`
		Arch         string `json:"arch"`
		NumCPU       int    `json:"num_cpu"`
		TimestampUTC string `json:"timestamp_utc"`
	}

	out := struct {
		Total      int64        `json:"total"`
		NextCursor string       `json:"next_cursor"`
		Items      []serverItem `json:"items"`
	}{Total: total, NextCursor: next, Items: make([]serverItem, 0, len(servers))}

	for i := range servers {
		out.Items = append(out.Items, serverItem{
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
