package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akaitigo/astro-karuta/backend/internal/handler"
	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
	"github.com/akaitigo/astro-karuta/backend/internal/seed"
	"github.com/akaitigo/astro-karuta/backend/internal/service"
)

func setupTestHandler(t *testing.T) (*http.ServeMux, *repository.InMemoryCardRepository) {
	t.Helper()
	cardRepo := repository.NewInMemoryCardRepository()
	deckRepo := repository.NewInMemoryDeckRepository()

	if err := seed.LoadCards(context.Background(), cardRepo); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	svc := service.NewCardService(cardRepo, deckRepo)
	h := handler.NewCardHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux, cardRepo
}

func TestListCards_All(t *testing.T) {
	mux, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cards", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var cards []model.Card
	if err := json.NewDecoder(rec.Body).Decode(&cards); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	// 30 constellations + 9 planets + 10 phenomena = 49
	if len(cards) < 49 {
		t.Errorf("expected at least 49 cards, got %d", len(cards))
	}
}

func TestListCards_FilterByCategory(t *testing.T) {
	mux, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cards?category=planet", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var cards []model.Card
	if err := json.NewDecoder(rec.Body).Decode(&cards); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(cards) != 9 {
		t.Errorf("expected 9 planets, got %d", len(cards))
	}
	for _, c := range cards {
		if c.Category != model.CardCategoryPlanet {
			t.Errorf("expected planet category, got %s", c.Category)
		}
	}
}

func TestListCards_FilterByCategory_Constellation(t *testing.T) {
	mux, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cards?category=constellation", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var cards []model.Card
	if err := json.NewDecoder(rec.Body).Decode(&cards); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(cards) != 30 {
		t.Errorf("expected 30 constellations, got %d", len(cards))
	}
}

func TestListCards_FilterByCategory_Phenomenon(t *testing.T) {
	mux, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cards?category=phenomenon", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var cards []model.Card
	if err := json.NewDecoder(rec.Body).Decode(&cards); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(cards) != 10 {
		t.Errorf("expected 10 phenomena, got %d", len(cards))
	}
}

func TestListCards_InvalidCategory(t *testing.T) {
	mux, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cards?category=invalid", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestListCards_FilterBySeason(t *testing.T) {
	mux, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cards?season=winter", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var cards []model.Card
	if err := json.NewDecoder(rec.Body).Decode(&cards); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(cards) == 0 {
		t.Error("expected at least some winter cards")
	}
	for _, c := range cards {
		if c.BestSeason != "winter" {
			t.Errorf("expected winter season, got %s for %s", c.BestSeason, c.Name)
		}
	}
}

func TestGetCard_Success(t *testing.T) {
	mux, cardRepo := setupTestHandler(t)

	// Get a card ID from the repo
	cards, _ := cardRepo.List(context.Background(), repository.CardFilter{})
	if len(cards) == 0 {
		t.Fatal("no cards seeded")
	}
	testCard := cards[0]

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cards/"+testCard.ID, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var card model.Card
	if err := json.NewDecoder(rec.Body).Decode(&card); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if card.ID != testCard.ID {
		t.Errorf("expected ID %s, got %s", testCard.ID, card.ID)
	}
}

func TestGetCard_NotFound(t *testing.T) {
	mux, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cards/nonexistent-id", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestListDecks_Empty(t *testing.T) {
	mux, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/decks", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var decks []model.Deck
	if err := json.NewDecoder(rec.Body).Decode(&decks); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(decks) != 0 {
		t.Errorf("expected 0 decks, got %d", len(decks))
	}
}
