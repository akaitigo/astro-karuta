package model

import "time"

// CollectionEntry represents a collected card by a user.
type CollectionEntry struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CardID     string    `json:"card_id"`
	ObtainedAt time.Time `json:"obtained_at"`
	Source     string    `json:"source"` // "game" or "mission"
}

// CollectionStats represents user's collection progress.
type CollectionStats struct {
	UserID      string  `json:"user_id"`
	TotalCards  int     `json:"total_cards"`
	Collected   int     `json:"collected"`
	Percentage  float64 `json:"percentage"`
}

// ObservationMission represents a real-sky observation challenge.
type ObservationMission struct {
	ID        string    `json:"id"`
	CardID    string    `json:"card_id"`
	Title     string    `json:"title"`
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}
