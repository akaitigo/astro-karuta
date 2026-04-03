package service_test

import (
	"context"
	"testing"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
	"github.com/akaitigo/astro-karuta/backend/internal/seed"
	"github.com/akaitigo/astro-karuta/backend/internal/service"
)

func setupCardService(t *testing.T) *service.CardService {
	t.Helper()
	cardRepo := repository.NewInMemoryCardRepository()
	deckRepo := repository.NewInMemoryDeckRepository()

	if err := seed.LoadCards(context.Background(), cardRepo); err != nil {
		t.Fatal(err)
	}
	return service.NewCardService(cardRepo, deckRepo)
}

func TestCardService_ListCards_All(t *testing.T) {
	svc := setupCardService(t)

	cards, err := svc.ListCards(context.Background(), "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) < 49 {
		t.Errorf("expected at least 49 cards, got %d", len(cards))
	}
}

func TestCardService_ListCards_InvalidCategory(t *testing.T) {
	svc := setupCardService(t)

	_, err := svc.ListCards(context.Background(), "invalid", "")
	if err == nil {
		t.Fatal("expected error for invalid category")
	}
}

func TestCardService_ListCards_InvalidSeason(t *testing.T) {
	svc := setupCardService(t)

	_, err := svc.ListCards(context.Background(), "", "invalid-season")
	if err == nil {
		t.Fatal("expected error for invalid season")
	}
}

func TestCardService_ListCards_FilterPlanets(t *testing.T) {
	svc := setupCardService(t)

	cards, err := svc.ListCards(context.Background(), "planet", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 9 {
		t.Errorf("expected 9 planets, got %d", len(cards))
	}
	for _, c := range cards {
		if c.Category != model.CardCategoryPlanet {
			t.Errorf("expected planet, got %s", c.Category)
		}
	}
}

func TestCardService_GetCard_EmptyID(t *testing.T) {
	svc := setupCardService(t)

	_, err := svc.GetCard(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestCardService_CardCount(t *testing.T) {
	svc := setupCardService(t)

	count, err := svc.CardCount(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if count < 49 {
		t.Errorf("expected at least 49 cards, got %d", count)
	}
}
