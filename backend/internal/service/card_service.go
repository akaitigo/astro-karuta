package service

import (
	"context"
	"fmt"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
)

// CardService provides business logic for card operations.
type CardService struct {
	cardRepo    repository.CardRepository
	deckRepo    repository.DeckRepository
	seasonalSvc *SeasonalService
}

// NewCardService creates a new CardService.
func NewCardService(cardRepo repository.CardRepository, deckRepo repository.DeckRepository) *CardService {
	return &CardService{
		cardRepo: cardRepo,
		deckRepo: deckRepo,
	}
}

// SetSeasonalService sets the seasonal service for deck generation.
func (s *CardService) SetSeasonalService(svc *SeasonalService) {
	s.seasonalSvc = svc
}

// ListCards returns cards filtered by the given criteria.
func (s *CardService) ListCards(ctx context.Context, category string, season string) ([]model.Card, error) {
	filter := repository.CardFilter{}
	if category != "" {
		cat := model.CardCategory(category)
		if !isValidCategory(cat) {
			return nil, fmt.Errorf("invalid category: %s", category)
		}
		filter.Category = cat
	}
	if season != "" {
		if !isValidSeason(season) {
			return nil, fmt.Errorf("invalid season: %s", season)
		}
		filter.BestSeason = season
	}
	return s.cardRepo.List(ctx, filter)
}

// GetCard returns a single card by ID.
func (s *CardService) GetCard(ctx context.Context, id string) (*model.Card, error) {
	if id == "" {
		return nil, fmt.Errorf("card ID must not be empty")
	}
	return s.cardRepo.GetByID(ctx, id)
}

// ListDecks returns all available decks.
func (s *CardService) ListDecks(ctx context.Context) ([]model.Deck, error) {
	return s.deckRepo.List(ctx)
}

// GetDeck returns a deck by ID.
func (s *CardService) GetDeck(ctx context.Context, id string) (*model.Deck, error) {
	if id == "" {
		return nil, fmt.Errorf("deck ID must not be empty")
	}
	return s.deckRepo.GetByID(ctx, id)
}

// GetSeasonalDeck returns the seasonal deck for the current month.
// If a SeasonalService is set, it delegates to that service for automatic
// deck generation based on visible constellations.
func (s *CardService) GetSeasonalDeck(ctx context.Context) (*model.Deck, error) {
	if s.seasonalSvc != nil {
		return s.seasonalSvc.GetSeasonalDeck(ctx)
	}
	month := time.Now().Month()
	return s.deckRepo.GetSeasonal(ctx, month)
}

// CardCount returns the total number of cards.
func (s *CardService) CardCount(ctx context.Context) (int, error) {
	return s.cardRepo.Count(ctx)
}

func isValidCategory(c model.CardCategory) bool {
	switch c {
	case model.CardCategoryConstellation, model.CardCategoryPlanet, model.CardCategoryPhenomenon:
		return true
	}
	return false
}

func isValidSeason(s string) bool {
	switch s {
	case "spring", "summer", "autumn", "winter", "all":
		return true
	}
	return false
}
