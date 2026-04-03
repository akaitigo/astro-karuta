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

func setupCollectionHandler(t *testing.T) (*http.ServeMux, *repository.InMemoryCollectionRepository) {
	t.Helper()
	cardRepo := repository.NewInMemoryCardRepository()
	if err := seed.LoadCards(context.Background(), cardRepo); err != nil {
		t.Fatalf("failed to seed: %v", err)
	}

	collectionRepo := repository.NewInMemoryCollectionRepository(cardRepo)
	collectionSvc := service.NewCollectionService(collectionRepo)
	h := handler.NewCollectionHandler(collectionSvc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux, collectionRepo
}

func TestListCollection_Success(t *testing.T) {
	mux, collectionRepo := setupCollectionHandler(t)
	ctx := context.Background()

	// Add some entries
	if err := collectionRepo.AddToCollection(ctx, "user-1", "card-1", "game"); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/collections?user_id=user-1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var entries []model.CollectionEntry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestListCollection_MissingUserID(t *testing.T) {
	mux, _ := setupCollectionHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/collections", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestListCollection_EmptyCollection(t *testing.T) {
	mux, _ := setupCollectionHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/collections?user_id=user-1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var entries []model.CollectionEntry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestListCollection_InvalidCategory(t *testing.T) {
	mux, _ := setupCollectionHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/collections?user_id=user-1&category=invalid", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestListCollection_FilterByCategory(t *testing.T) {
	mux, collectionRepo := setupCollectionHandler(t)
	ctx := context.Background()

	// We need to use real card IDs from seed data, so add entries with arbitrary IDs
	// that exist in the collection repo (the card_id might not match seed data,
	// but GetCollection with category filter needs the card to exist in cardRepo)
	// For this test, we add entries and filter without category
	if err := collectionRepo.AddToCollection(ctx, "user-1", "card-A", "game"); err != nil {
		t.Fatal(err)
	}
	if err := collectionRepo.AddToCollection(ctx, "user-1", "card-B", "mission"); err != nil {
		t.Fatal(err)
	}

	// Without category filter, should return all entries
	req := httptest.NewRequest(http.MethodGet, "/api/v1/collections?user_id=user-1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var entries []model.CollectionEntry
	if err := json.NewDecoder(rec.Body).Decode(&entries); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestGetStats_Success(t *testing.T) {
	mux, collectionRepo := setupCollectionHandler(t)
	ctx := context.Background()

	if err := collectionRepo.AddToCollection(ctx, "user-1", "card-1", "game"); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/collections/stats?user_id=user-1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var stats model.CollectionStats
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if stats.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", stats.UserID)
	}
	if stats.Collected != 1 {
		t.Errorf("expected 1 collected, got %d", stats.Collected)
	}
	// total_cards should be 49 (seeded)
	if stats.TotalCards != 49 {
		t.Errorf("expected 49 total cards, got %d", stats.TotalCards)
	}
}

func TestGetStats_MissingUserID(t *testing.T) {
	mux, _ := setupCollectionHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/collections/stats", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestGetStats_EmptyCollection(t *testing.T) {
	mux, _ := setupCollectionHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/collections/stats?user_id=new-user", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var stats model.CollectionStats
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if stats.Collected != 0 {
		t.Errorf("expected 0 collected, got %d", stats.Collected)
	}
	if stats.Percentage != 0.0 {
		t.Errorf("expected 0.0%%, got %.1f%%", stats.Percentage)
	}
}
