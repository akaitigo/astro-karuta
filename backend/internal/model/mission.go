package model

import "time"

// MissionStatus represents the state of a user's mission.
type MissionStatus string

const (
	MissionStatusActive    MissionStatus = "active"
	MissionStatusCompleted MissionStatus = "completed"
	MissionStatusExpired   MissionStatus = "expired"
)

// UserMission represents an observation mission assigned to a user.
type UserMission struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	MissionID   string        `json:"mission_id"`
	CardID      string        `json:"card_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      MissionStatus `json:"status"`
	ValidFrom   time.Time     `json:"valid_from"`
	ValidTo     time.Time     `json:"valid_to"`
	CompletedAt *time.Time    `json:"completed_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}

// CompleteMissionRequest is the request body for completing a mission.
type CompleteMissionRequest struct {
	UserID    string  `json:"user_id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

// CompleteMissionResponse is returned after successfully completing a mission.
type CompleteMissionResponse struct {
	Mission   UserMission `json:"mission"`
	BonusCard *Card       `json:"bonus_card,omitempty"`
}
