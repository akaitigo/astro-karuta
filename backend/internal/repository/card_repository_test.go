package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
)

func TestInMemoryCardRepository_Create(t *testing.T) {
	repo := repository.NewInMemoryCardRepository()
	ctx := context.Background()

	card := &model.Card{
		ID:          "card-1",
		Name:        "オリオン座",
		Category:    model.CardCategoryConstellation,
		ReadingText: "冬の夜空に三つ星が並ぶ",
		BestSeason:  "winter",
		CreatedAt:   time.Now(),
	}

	err := repo.Create(ctx, card)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}

func TestInMemoryCardRepository_Create_EmptyID(t *testing.T) {
	repo := repository.NewInMemoryCardRepository()
	ctx := context.Background()

	card := &model.Card{Name: "Test"}
	err := repo.Create(ctx, card)
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestInMemoryCardRepository_Create_Duplicate(t *testing.T) {
	repo := repository.NewInMemoryCardRepository()
	ctx := context.Background()

	card := &model.Card{ID: "card-1", Name: "Test", CreatedAt: time.Now()}
	if err := repo.Create(ctx, card); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := repo.Create(ctx, card)
	if err == nil {
		t.Fatal("expected error for duplicate card")
	}
}

func TestInMemoryCardRepository_GetByID(t *testing.T) {
	repo := repository.NewInMemoryCardRepository()
	ctx := context.Background()

	card := &model.Card{
		ID:       "card-1",
		Name:     "オリオン座",
		Category: model.CardCategoryConstellation,
	}
	if err := repo.Create(ctx, card); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetByID(ctx, "card-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "オリオン座" {
		t.Errorf("expected name オリオン座, got %s", got.Name)
	}
}

func TestInMemoryCardRepository_GetByID_NotFound(t *testing.T) {
	repo := repository.NewInMemoryCardRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent card")
	}
}

func TestInMemoryCardRepository_List(t *testing.T) {
	repo := repository.NewInMemoryCardRepository()
	ctx := context.Background()

	cards := []*model.Card{
		{ID: "c1", Name: "オリオン座", Category: model.CardCategoryConstellation, BestSeason: "winter"},
		{ID: "c2", Name: "火星", Category: model.CardCategoryPlanet, BestSeason: "all"},
		{ID: "c3", Name: "さそり座", Category: model.CardCategoryConstellation, BestSeason: "summer"},
	}
	for _, c := range cards {
		if err := repo.Create(ctx, c); err != nil {
			t.Fatal(err)
		}
	}

	// No filter
	all, err := repo.List(ctx, repository.CardFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 cards, got %d", len(all))
	}

	// Filter by category
	constellations, err := repo.List(ctx, repository.CardFilter{Category: model.CardCategoryConstellation})
	if err != nil {
		t.Fatal(err)
	}
	if len(constellations) != 2 {
		t.Errorf("expected 2 constellations, got %d", len(constellations))
	}

	// Filter by season
	winter, err := repo.List(ctx, repository.CardFilter{BestSeason: "winter"})
	if err != nil {
		t.Fatal(err)
	}
	if len(winter) != 1 {
		t.Errorf("expected 1 winter card, got %d", len(winter))
	}
}
