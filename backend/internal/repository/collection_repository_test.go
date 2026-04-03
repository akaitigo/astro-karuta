package repository_test

import (
	"context"
	"testing"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
)

func setupCollectionTest(t *testing.T) (*repository.InMemoryCollectionRepository, *repository.InMemoryCardRepository) {
	t.Helper()
	cardRepo := repository.NewInMemoryCardRepository()
	ctx := context.Background()

	cards := []*model.Card{
		{ID: "card-1", Name: "オリオン座", Category: model.CardCategoryConstellation, BestSeason: "winter"},
		{ID: "card-2", Name: "火星", Category: model.CardCategoryPlanet, BestSeason: "all"},
		{ID: "card-3", Name: "皆既日食", Category: model.CardCategoryPhenomenon, BestSeason: "all"},
		{ID: "card-4", Name: "さそり座", Category: model.CardCategoryConstellation, BestSeason: "summer"},
	}
	for _, c := range cards {
		if err := cardRepo.Create(ctx, c); err != nil {
			t.Fatalf("failed to seed card: %v", err)
		}
	}

	collectionRepo := repository.NewInMemoryCollectionRepository(cardRepo)
	return collectionRepo, cardRepo
}

func TestAddToCollection_Success(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	err := repo.AddToCollection(ctx, "user-1", "card-1", "game")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := repo.GetCollection(ctx, "user-1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].CardID != "card-1" {
		t.Errorf("expected card-1, got %s", entries[0].CardID)
	}
	if entries[0].Source != "game" {
		t.Errorf("expected source game, got %s", entries[0].Source)
	}
}

func TestAddToCollection_EmptyUserID(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	err := repo.AddToCollection(ctx, "", "card-1", "game")
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}

func TestAddToCollection_EmptyCardID(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	err := repo.AddToCollection(ctx, "user-1", "", "game")
	if err == nil {
		t.Fatal("expected error for empty card_id")
	}
}

func TestAddToCollection_EmptySource(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	err := repo.AddToCollection(ctx, "user-1", "card-1", "")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestAddToCollection_DuplicateSkip(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	if err := repo.AddToCollection(ctx, "user-1", "card-1", "game"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Adding the same card again should not error and should not duplicate
	if err := repo.AddToCollection(ctx, "user-1", "card-1", "mission"); err != nil {
		t.Fatalf("unexpected error on duplicate: %v", err)
	}

	entries, err := repo.GetCollection(ctx, "user-1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry after duplicate add, got %d", len(entries))
	}
}

func TestGetCollection_FilterByCategory(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	// Add cards from different categories
	if err := repo.AddToCollection(ctx, "user-1", "card-1", "game"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AddToCollection(ctx, "user-1", "card-2", "game"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AddToCollection(ctx, "user-1", "card-3", "mission"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AddToCollection(ctx, "user-1", "card-4", "game"); err != nil {
		t.Fatal(err)
	}

	// Filter constellations only
	constellations, err := repo.GetCollection(ctx, "user-1", model.CardCategoryConstellation)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(constellations) != 2 {
		t.Errorf("expected 2 constellations, got %d", len(constellations))
	}

	// Filter planets only
	planets, err := repo.GetCollection(ctx, "user-1", model.CardCategoryPlanet)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(planets) != 1 {
		t.Errorf("expected 1 planet, got %d", len(planets))
	}

	// No filter — all entries
	all, err := repo.GetCollection(ctx, "user-1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 4 {
		t.Errorf("expected 4 entries, got %d", len(all))
	}
}

func TestGetCollection_EmptyUserID(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	_, err := repo.GetCollection(ctx, "", "")
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}

func TestGetCollection_DifferentUsers(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	if err := repo.AddToCollection(ctx, "user-1", "card-1", "game"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AddToCollection(ctx, "user-2", "card-2", "mission"); err != nil {
		t.Fatal(err)
	}

	user1Entries, err := repo.GetCollection(ctx, "user-1", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(user1Entries) != 1 {
		t.Errorf("user-1: expected 1, got %d", len(user1Entries))
	}

	user2Entries, err := repo.GetCollection(ctx, "user-2", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(user2Entries) != 1 {
		t.Errorf("user-2: expected 1, got %d", len(user2Entries))
	}
}

func TestGetStats_Success(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	if err := repo.AddToCollection(ctx, "user-1", "card-1", "game"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AddToCollection(ctx, "user-1", "card-2", "game"); err != nil {
		t.Fatal(err)
	}

	stats, err := repo.GetStats(ctx, "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", stats.UserID)
	}
	if stats.TotalCards != 4 {
		t.Errorf("expected 4 total cards, got %d", stats.TotalCards)
	}
	if stats.Collected != 2 {
		t.Errorf("expected 2 collected, got %d", stats.Collected)
	}
	if stats.Percentage != 50.0 {
		t.Errorf("expected 50.0%%, got %.1f%%", stats.Percentage)
	}
}

func TestGetStats_EmptyCollection(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	stats, err := repo.GetStats(ctx, "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Collected != 0 {
		t.Errorf("expected 0 collected, got %d", stats.Collected)
	}
	if stats.Percentage != 0.0 {
		t.Errorf("expected 0.0%%, got %.1f%%", stats.Percentage)
	}
}

func TestGetStats_EmptyUserID(t *testing.T) {
	repo, _ := setupCollectionTest(t)
	ctx := context.Background()

	_, err := repo.GetStats(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}
