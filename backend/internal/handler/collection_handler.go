package handler

import (
	"net/http"

	"github.com/akaitigo/astro-karuta/backend/internal/service"
)

// CollectionHandler handles HTTP requests for collection endpoints.
type CollectionHandler struct {
	svc *service.CollectionService
}

// NewCollectionHandler creates a new CollectionHandler.
func NewCollectionHandler(svc *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{svc: svc}
}

// RegisterRoutes registers collection routes to the mux.
func (h *CollectionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/collections", h.ListCollection)
	mux.HandleFunc("GET /api/v1/collections/stats", h.GetStats)
}

// ListCollection returns the user's collection entries.
func (h *CollectionHandler) ListCollection(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	category := r.URL.Query().Get("category")

	entries, err := h.svc.GetCollection(r.Context(), userID, category)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, entries)
}

// GetStats returns collection statistics for a user.
func (h *CollectionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	stats, err := h.svc.GetStats(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
