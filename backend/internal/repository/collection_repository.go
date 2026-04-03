package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/google/uuid"
)

// CollectionRepository defines operations for collection persistence.
type CollectionRepository interface {
	AddToCollection(ctx context.Context, userID, cardID, source string) error
	GetCollection(ctx context.Context, userID string, category model.CardCategory) ([]model.CollectionEntry, error)
	GetStats(ctx context.Context, userID string) (*model.CollectionStats, error)
}

// InMemoryCollectionRepository is an in-memory implementation of CollectionRepository.
type InMemoryCollectionRepository struct {
	mu       sync.RWMutex
	entries  []model.CollectionEntry
	cardRepo CardRepository
}

// NewInMemoryCollectionRepository creates a new in-memory collection repository.
func NewInMemoryCollectionRepository(cardRepo CardRepository) *InMemoryCollectionRepository {
	return &InMemoryCollectionRepository{
		entries:  make([]model.CollectionEntry, 0),
		cardRepo: cardRepo,
	}
}

// AddToCollection adds a card to a user's collection. Skips if already collected.
func (r *InMemoryCollectionRepository) AddToCollection(_ context.Context, userID, cardID, source string) error {
	if userID == "" {
		return fmt.Errorf("user_id must not be empty")
	}
	if cardID == "" {
		return fmt.Errorf("card_id must not be empty")
	}
	if source == "" {
		return fmt.Errorf("source must not be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Skip if duplicate
	for _, e := range r.entries {
		if e.UserID == userID && e.CardID == cardID {
			return nil
		}
	}

	entry := model.CollectionEntry{
		ID:         uuid.New().String(),
		UserID:     userID,
		CardID:     cardID,
		ObtainedAt: time.Now(),
		Source:     source,
	}
	r.entries = append(r.entries, entry)
	return nil
}

// GetCollection returns collection entries for a user, optionally filtered by card category.
func (r *InMemoryCollectionRepository) GetCollection(ctx context.Context, userID string, category model.CardCategory) ([]model.CollectionEntry, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id must not be empty")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.CollectionEntry, 0)
	for _, e := range r.entries {
		if e.UserID != userID {
			continue
		}
		if category != "" {
			card, err := r.cardRepo.GetByID(ctx, e.CardID)
			if err != nil {
				continue
			}
			if card.Category != category {
				continue
			}
		}
		result = append(result, e)
	}
	return result, nil
}

// GetStats returns collection statistics for a user.
func (r *InMemoryCollectionRepository) GetStats(ctx context.Context, userID string) (*model.CollectionStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id must not be empty")
	}

	totalCards, err := r.cardRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count total cards: %w", err)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	collected := 0
	for _, e := range r.entries {
		if e.UserID == userID {
			collected++
		}
	}

	var percentage float64
	if totalCards > 0 {
		percentage = float64(collected) / float64(totalCards) * 100
	}

	return &model.CollectionStats{
		UserID:     userID,
		TotalCards: totalCards,
		Collected:  collected,
		Percentage: percentage,
	}, nil
}
