package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
)

// CardFilter defines filter criteria for card queries.
type CardFilter struct {
	Category   model.CardCategory
	BestSeason string
}

// CardRepository defines operations for card persistence.
type CardRepository interface {
	List(ctx context.Context, filter CardFilter) ([]model.Card, error)
	GetByID(ctx context.Context, id string) (*model.Card, error)
	Create(ctx context.Context, card *model.Card) error
	Count(ctx context.Context) (int, error)
}

// InMemoryCardRepository is an in-memory implementation of CardRepository.
type InMemoryCardRepository struct {
	mu    sync.RWMutex
	cards map[string]model.Card
	order []string
}

// NewInMemoryCardRepository creates a new in-memory card repository.
func NewInMemoryCardRepository() *InMemoryCardRepository {
	return &InMemoryCardRepository{
		cards: make(map[string]model.Card),
	}
}

func (r *InMemoryCardRepository) List(_ context.Context, filter CardFilter) ([]model.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Card, 0, len(r.cards))
	for _, id := range r.order {
		c := r.cards[id]
		if filter.Category != "" && c.Category != filter.Category {
			continue
		}
		if filter.BestSeason != "" && !strings.EqualFold(c.BestSeason, filter.BestSeason) {
			continue
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *InMemoryCardRepository) GetByID(_ context.Context, id string) (*model.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	card, ok := r.cards[id]
	if !ok {
		return nil, fmt.Errorf("card not found: %s", id)
	}
	return &card, nil
}

func (r *InMemoryCardRepository) Create(_ context.Context, card *model.Card) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if card.ID == "" {
		return fmt.Errorf("card ID must not be empty")
	}
	if _, exists := r.cards[card.ID]; exists {
		return fmt.Errorf("card already exists: %s", card.ID)
	}
	r.cards[card.ID] = *card
	r.order = append(r.order, card.ID)
	return nil
}

func (r *InMemoryCardRepository) Count(_ context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.cards), nil
}
