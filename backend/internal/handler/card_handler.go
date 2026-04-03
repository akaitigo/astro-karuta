package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/astro-karuta/backend/internal/service"
)

// CardHandler handles HTTP requests for card and deck endpoints.
type CardHandler struct {
	svc *service.CardService
}

// NewCardHandler creates a new CardHandler.
func NewCardHandler(svc *service.CardService) *CardHandler {
	return &CardHandler{svc: svc}
}

// RegisterRoutes registers card and deck routes to the mux.
func (h *CardHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/cards", h.ListCards)
	mux.HandleFunc("GET /api/v1/cards/{id}", h.GetCard)
	mux.HandleFunc("GET /api/v1/decks", h.ListDecks)
	mux.HandleFunc("GET /api/v1/decks/seasonal", h.GetSeasonalDeck)
	mux.HandleFunc("GET /api/v1/decks/{id}", h.GetDeck)
}

func (h *CardHandler) ListCards(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	season := r.URL.Query().Get("season")

	cards, err := h.svc.ListCards(r.Context(), category, season)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cards)
}

func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	card, err := h.svc.GetCard(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, card)
}

func (h *CardHandler) ListDecks(w http.ResponseWriter, r *http.Request) {
	decks, err := h.svc.ListDecks(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, decks)
}

func (h *CardHandler) GetDeck(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	deck, err := h.svc.GetDeck(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, deck)
}

func (h *CardHandler) GetSeasonalDeck(w http.ResponseWriter, r *http.Request) {
	deck, err := h.svc.GetSeasonalDeck(r.Context())
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, deck)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode error response", http.StatusInternalServerError)
	}
}
