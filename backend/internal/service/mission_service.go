package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
	"github.com/akaitigo/astro-karuta/backend/pkg/astronomy"
	"github.com/google/uuid"
)

// maxLatitudeDeviation is the maximum allowed deviation in latitude degrees
// for location validation when completing a mission.
const maxLatitudeDeviation = 30.0

// MissionService provides business logic for observation missions.
type MissionService struct {
	missionRepo repository.MissionRepository
	cardRepo    repository.CardRepository
	nowFunc     func() time.Time
}

// NewMissionService creates a new MissionService.
func NewMissionService(missionRepo repository.MissionRepository, cardRepo repository.CardRepository) *MissionService {
	return &MissionService{
		missionRepo: missionRepo,
		cardRepo:    cardRepo,
		nowFunc:     time.Now,
	}
}

// SetNowFunc overrides the time source for testing.
func (s *MissionService) SetNowFunc(fn func() time.Time) {
	s.nowFunc = fn
}

// GetActiveMissions returns active missions for a user.
// If the user has no missions, it generates missions based on the current month's
// visible constellations.
func (s *MissionService) GetActiveMissions(ctx context.Context, userID string) ([]model.UserMission, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID must not be empty")
	}

	missions, err := s.missionRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list active missions: %w", err)
	}

	// If user has active missions, return them
	if len(missions) > 0 {
		return missions, nil
	}

	// Generate new missions based on current month's visible constellations
	now := s.nowFunc()
	month := int(now.Month())
	visible := astronomy.GetVisibleConstellations(month, astronomy.DefaultLatitude)
	if len(visible) == 0 {
		return []model.UserMission{}, nil
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

	// Create up to 5 missions from visible constellations
	validFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	validTo := validFrom.AddDate(0, 1, 0).Add(-time.Second)

	var generated []model.UserMission
	count := 0
	for _, card := range allCards {
		if count >= 5 {
			break
		}
		if !visibleSet[card.Name] {
			continue
		}
		mission := model.UserMission{
			ID:          uuid.New().String(),
			UserID:      userID,
			MissionID:   fmt.Sprintf("obs-%d-%s", month, card.ID),
			CardID:      card.ID,
			Title:       fmt.Sprintf("%sを観測しよう", card.Name),
			Description: fmt.Sprintf("夜空で%sを見つけて、観測ボタンを押そう！", card.Name),
			Status:      model.MissionStatusActive,
			ValidFrom:   validFrom,
			ValidTo:     validTo,
			CreatedAt:   now,
		}
		if err := s.missionRepo.Create(ctx, &mission); err != nil {
			return nil, fmt.Errorf("create mission: %w", err)
		}
		generated = append(generated, mission)
		count++
	}

	return generated, nil
}

// CompleteMission validates location and time, marks the mission as completed,
// and returns the bonus card.
func (s *MissionService) CompleteMission(ctx context.Context, userID string, missionID string, lat float64, lng float64) (*model.CompleteMissionResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID must not be empty")
	}
	if missionID == "" {
		return nil, fmt.Errorf("mission ID must not be empty")
	}

	mission, err := s.missionRepo.GetByID(ctx, missionID)
	if err != nil {
		return nil, fmt.Errorf("get mission: %w", err)
	}

	// Verify the mission belongs to this user
	if mission.UserID != userID {
		return nil, fmt.Errorf("mission does not belong to user")
	}

	// Verify mission is still active
	if mission.Status != model.MissionStatusActive {
		return nil, fmt.Errorf("mission is not active (status: %s)", mission.Status)
	}

	now := s.nowFunc()

	// Verify mission hasn't expired
	if now.After(mission.ValidTo) {
		return nil, fmt.Errorf("mission has expired")
	}

	// Time validation: must be nighttime (18:00-06:00)
	if err := validateNighttime(now); err != nil {
		return nil, err
	}

	// Location validation: latitude within +-30 degrees of default
	if err := validateLocation(lat, lng); err != nil {
		return nil, err
	}

	// Mark mission as completed
	completedAt := now
	mission.Status = model.MissionStatusCompleted
	mission.CompletedAt = &completedAt

	if err := s.missionRepo.Update(ctx, mission); err != nil {
		return nil, fmt.Errorf("update mission: %w", err)
	}

	// Get the bonus card
	card, err := s.cardRepo.GetByID(ctx, mission.CardID)
	if err != nil {
		return nil, fmt.Errorf("get bonus card: %w", err)
	}

	return &model.CompleteMissionResponse{
		Mission:   *mission,
		BonusCard: card,
	}, nil
}

// validateNighttime checks that the time is between 18:00 and 06:00.
func validateNighttime(t time.Time) error {
	hour := t.Hour()
	if hour >= 6 && hour < 18 {
		return fmt.Errorf("observation must be at nighttime (18:00-06:00), current hour: %d", hour)
	}
	return nil
}

// validateLocation checks that the latitude is within +-30 degrees
// of the default observation latitude.
func validateLocation(lat float64, _ float64) error {
	deviation := math.Abs(lat - astronomy.DefaultLatitude)
	if deviation > maxLatitudeDeviation {
		return fmt.Errorf("location too far from observation area (deviation: %.1f degrees, max: %.1f)", deviation, maxLatitudeDeviation)
	}
	return nil
}
