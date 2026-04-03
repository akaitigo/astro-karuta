package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
)

// DeckRepository defines operations for deck persistence.
type DeckRepository interface {
	List(ctx context.Context) ([]model.Deck, error)
	GetByID(ctx context.Context, id string) (*model.Deck, error)
	GetSeasonal(ctx context.Context, month time.Month) (*model.Deck, error)
	Create(ctx context.Context, deck *model.Deck) error
}

// InMemoryDeckRepository is an in-memory implementation of DeckRepository.
type InMemoryDeckRepository struct {
	mu    sync.RWMutex
	decks map[string]model.Deck
	order []string
}

// NewInMemoryDeckRepository creates a new in-memory deck repository.
func NewInMemoryDeckRepository() *InMemoryDeckRepository {
	return &InMemoryDeckRepository{
		decks: make(map[string]model.Deck),
	}
}

func (r *InMemoryDeckRepository) List(_ context.Context) ([]model.Deck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Deck, 0, len(r.decks))
	for _, id := range r.order {
		result = append(result, r.decks[id])
	}
	return result, nil
}

func (r *InMemoryDeckRepository) GetByID(_ context.Context, id string) (*model.Deck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	deck, ok := r.decks[id]
	if !ok {
		return nil, fmt.Errorf("deck not found: %s", id)
	}
	return &deck, nil
}

func (r *InMemoryDeckRepository) GetSeasonal(_ context.Context, month time.Month) (*model.Deck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()
	for _, deck := range r.decks {
		if deck.Seasonal && !now.Before(deck.ValidFrom) && !now.After(deck.ValidTo) {
			if deck.ValidFrom.Month() == month || deck.ValidTo.Month() == month {
				return &deck, nil
			}
		}
	}
	return nil, fmt.Errorf("no seasonal deck for month: %s", month)
}

func (r *InMemoryDeckRepository) Create(_ context.Context, deck *model.Deck) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if deck.ID == "" {
		return fmt.Errorf("deck ID must not be empty")
	}
	if _, exists := r.decks[deck.ID]; exists {
		return fmt.Errorf("deck already exists: %s", deck.ID)
	}
	r.decks[deck.ID] = *deck
	r.order = append(r.order, deck.ID)
	return nil
}
