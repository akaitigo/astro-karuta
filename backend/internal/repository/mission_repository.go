package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
)

// MissionRepository defines operations for mission persistence.
type MissionRepository interface {
	Create(ctx context.Context, mission *model.UserMission) error
	GetByID(ctx context.Context, id string) (*model.UserMission, error)
	ListByUserID(ctx context.Context, userID string) ([]model.UserMission, error)
	ListActiveByUserID(ctx context.Context, userID string) ([]model.UserMission, error)
	Update(ctx context.Context, mission *model.UserMission) error
}

// InMemoryMissionRepository is an in-memory implementation of MissionRepository.
type InMemoryMissionRepository struct {
	mu       sync.RWMutex
	missions map[string]model.UserMission
	order    []string
}

// NewInMemoryMissionRepository creates a new in-memory mission repository.
func NewInMemoryMissionRepository() *InMemoryMissionRepository {
	return &InMemoryMissionRepository{
		missions: make(map[string]model.UserMission),
	}
}

func (r *InMemoryMissionRepository) Create(_ context.Context, mission *model.UserMission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if mission.ID == "" {
		return fmt.Errorf("mission ID must not be empty")
	}
	if _, exists := r.missions[mission.ID]; exists {
		return fmt.Errorf("mission already exists: %s", mission.ID)
	}
	r.missions[mission.ID] = *mission
	r.order = append(r.order, mission.ID)
	return nil
}

func (r *InMemoryMissionRepository) GetByID(_ context.Context, id string) (*model.UserMission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mission, ok := r.missions[id]
	if !ok {
		return nil, fmt.Errorf("mission not found: %s", id)
	}
	return &mission, nil
}

func (r *InMemoryMissionRepository) ListByUserID(_ context.Context, userID string) ([]model.UserMission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.UserMission
	for _, id := range r.order {
		m := r.missions[id]
		if m.UserID == userID {
			result = append(result, m)
		}
	}
	return result, nil
}

func (r *InMemoryMissionRepository) ListActiveByUserID(_ context.Context, userID string) ([]model.UserMission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()
	var result []model.UserMission
	for _, id := range r.order {
		m := r.missions[id]
		if m.UserID == userID && m.Status == model.MissionStatusActive && !now.After(m.ValidTo) {
			result = append(result, m)
		}
	}
	return result, nil
}

func (r *InMemoryMissionRepository) Update(_ context.Context, mission *model.UserMission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.missions[mission.ID]; !exists {
		return fmt.Errorf("mission not found: %s", mission.ID)
	}
	r.missions[mission.ID] = *mission
	return nil
}
