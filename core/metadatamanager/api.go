package metadatamanager

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service MetadataService
}

func NewHandler(s MetadataService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) BindRoutes(r chi.Router) {
	r.Post("/song/{id}/tag", h.UpdateSong)
	r.Post("/song/{id}/artwork", h.UpdateArtwork)
	r.Post("/youtube/download", h.DownloadYouTube)
}

type YTDownloadRequest struct {
	URL     string `json:"url"`
	Format  string `json:"format"`  // ex: "mp3", "flac", "m4a"
	Quality string `json:"quality"` // ex: "0" (better), "5" (mid)
}

func (h *Handler) DownloadYouTube(w http.ResponseWriter, r *http.Request) {
	var req YTDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "Missing YouTube URL", http.StatusBadRequest)
		return
	}

	go func() {
		ctx := r.Context()
		if err := h.service.DownloadFromYouTube(ctx, req.URL, req.Format, req.Quality); err != nil {
			return
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "Download started in background"}`))
}

func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		http.Error(w, "Missing song identifier", http.StatusBadRequest)
		return
	}

	var tags map[string]string
	if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTags(r.Context(), songID, tags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateArtwork(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	if songID == "" {
		http.Error(w, "Missing song identifier", http.StatusBadRequest)
		return
	}

	mimeType := r.Header.Get("Content-Type")
	if err := h.service.UpdateArtwork(r.Context(), songID, r.Body, mimeType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusNoContent)
}
