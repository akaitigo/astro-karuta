package service

import (
	"context"
	"fmt"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
)

// CollectionService provides business logic for collection operations.
type CollectionService struct {
	collectionRepo repository.CollectionRepository
}

// NewCollectionService creates a new CollectionService.
func NewCollectionService(collectionRepo repository.CollectionRepository) *CollectionService {
	return &CollectionService{
		collectionRepo: collectionRepo,
	}
}

// AddToCollection adds a card to a user's collection.
func (s *CollectionService) AddToCollection(ctx context.Context, userID, cardID, source string) error {
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}
	if cardID == "" {
		return fmt.Errorf("card_id is required")
	}
	if source != "game" && source != "mission" {
		return fmt.Errorf("source must be 'game' or 'mission', got: %s", source)
	}
	return s.collectionRepo.AddToCollection(ctx, userID, cardID, source)
}

// GetCollection returns a user's collection entries, optionally filtered by category.
func (s *CollectionService) GetCollection(ctx context.Context, userID string, category string) ([]model.CollectionEntry, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	var cat model.CardCategory
	if category != "" {
		cat = model.CardCategory(category)
		if !isValidCategory(cat) {
			return nil, fmt.Errorf("invalid category: %s", category)
		}
	}

	return s.collectionRepo.GetCollection(ctx, userID, cat)
}

// GetStats returns collection statistics for a user.
func (s *CollectionService) GetStats(ctx context.Context, userID string) (*model.CollectionStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	return s.collectionRepo.GetStats(ctx, userID)
}
