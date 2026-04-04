package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/service"
	"github.com/google/uuid"
)

// SeasonalHandler handles HTTP requests for seasonal and mission endpoints.
type SeasonalHandler struct {
	missionSvc  *service.MissionService
	seasonalSvc *service.SeasonalService
}

// NewSeasonalHandler creates a new SeasonalHandler.
func NewSeasonalHandler(missionSvc *service.MissionService, seasonalSvc *service.SeasonalService) *SeasonalHandler {
	return &SeasonalHandler{
		missionSvc:  missionSvc,
		seasonalSvc: seasonalSvc,
	}
}

// RegisterRoutes registers seasonal and mission routes to the mux.
func (h *SeasonalHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/missions", h.ListMissions)
	mux.HandleFunc("POST /api/v1/missions/{id}/complete", h.CompleteMission)
}

// ListMissions returns active missions for a user.
func (h *SeasonalHandler) ListMissions(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}
	// C5: validate user_id is a valid UUID
	if _, err := uuid.Parse(userID); err != nil {
		writeError(w, http.StatusBadRequest, "user_id must be a valid UUID")
		return
	}

	missions, err := h.missionSvc.GetActiveMissions(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, missions)
}

// CompleteMission validates and completes a mission.
func (h *SeasonalHandler) CompleteMission(w http.ResponseWriter, r *http.Request) {
	missionID := r.PathValue("id")
	if missionID == "" {
		writeError(w, http.StatusBadRequest, "mission ID is required")
		return
	}

	var req model.CompleteMissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}
	// C5: validate user_id is a valid UUID
	if _, err := uuid.Parse(req.UserID); err != nil {
		writeError(w, http.StatusBadRequest, "user_id must be a valid UUID")
		return
	}

	resp, err := h.missionSvc.CompleteMission(r.Context(), req.UserID, missionID, req.Latitude, req.Longitude)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
