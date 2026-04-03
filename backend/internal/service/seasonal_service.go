package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
	"github.com/akaitigo/astro-karuta/backend/pkg/astronomy"
	"github.com/google/uuid"
)

// SeasonalService provides seasonal deck generation logic.
type SeasonalService struct {
	cardRepo repository.CardRepository
	deckRepo repository.DeckRepository
}

// NewSeasonalService creates a new SeasonalService.
func NewSeasonalService(cardRepo repository.CardRepository, deckRepo repository.DeckRepository) *SeasonalService {
	return &SeasonalService{
		cardRepo: cardRepo,
		deckRepo: deckRepo,
	}
}

// GenerateSeasonalDeck creates or retrieves a seasonal deck
// containing cards for constellations visible in the given month.
func (s *SeasonalService) GenerateSeasonalDeck(ctx context.Context, month int) (*model.Deck, error) {
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month: %d", month)
	}

	// Try to get an existing seasonal deck for this month
	m := time.Month(month)
	existing, err := s.deckRepo.GetSeasonal(ctx, m)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Get visible constellation names
	visible := astronomy.GetVisibleConstellations(month, astronomy.DefaultLatitude)
	if len(visible) == 0 {
		return nil, fmt.Errorf("no visible constellations for month %d", month)
	}

	// Find matching constellation cards
	allCards, err := s.cardRepo.List(ctx, repository.CardFilter{
		Category: model.CardCategoryConstellation,
	})
	if err != nil {
		return nil, fmt.Errorf("list constellation cards: %w", err)
	}

	visibleSet := make(map[string]bool, len(visible))
	for _, name := range visible {
		visibleSet[name] = true
	}

	var cardIDs []string
	for _, card := range allCards {
		if visibleSet[strings.TrimSpace(card.Name)] {
			cardIDs = append(cardIDs, card.ID)
		}
	}

	if len(cardIDs) == 0 {
		return nil, fmt.Errorf("no matching cards found for month %d", month)
	}

	// Build the valid time range for this month
	now := time.Now()
	year := now.Year()
	validFrom := time.Date(year, m, 1, 0, 0, 0, 0, time.Local)
	validTo := validFrom.AddDate(0, 1, 0).Add(-time.Second)

	deck := &model.Deck{
		ID:        uuid.New().String(),
		Name:      fmt.Sprintf("%d月の星座デッキ", month),
		CardIDs:   cardIDs,
		Seasonal:  true,
		ValidFrom: validFrom,
		ValidTo:   validTo,
	}

	if err := s.deckRepo.Create(ctx, deck); err != nil {
		return nil, fmt.Errorf("create seasonal deck: %w", err)
	}

	return deck, nil
}

// GetSeasonalDeck returns the seasonal deck for the current month,
// generating it if it does not exist.
func (s *SeasonalService) GetSeasonalDeck(ctx context.Context) (*model.Deck, error) {
	month := int(time.Now().Month())
	return s.GenerateSeasonalDeck(ctx, month)
}
