package metadatamanager

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service MetadataService
}

func NewHandler(s MetadataService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) BindRoutes(mux *http.ServeMux) {
	// Use a clear prefix to avoid conflicts with Subsonic routes
	mux.HandleFunc("POST /api/extra/metadata/{id}", h.UpdateSong)
}

func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	songID := r.PathValue("id")
	if songID == "" {
		http.Error(w, "Missing song identifier", http.StatusBadRequest)
		return
	}

	var tags map[string]string
	if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTags(songID, tags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}