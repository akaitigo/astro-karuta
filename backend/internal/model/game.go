package model

import "time"

// GameStatus represents the current state of a game.
type GameStatus string

const (
	GameStatusWaiting  GameStatus = "waiting"
	GameStatusPlaying  GameStatus = "playing"
	GameStatusFinished GameStatus = "finished"
)

// Game represents an active karuta game session.
type Game struct {
	ID            string     `json:"id"`
	RoomCode      string     `json:"room_code"`
	Status        GameStatus `json:"status"`
	DeckID        string     `json:"deck_id"`
	Players       []Player   `json:"players"`
	CurrentCardID string     `json:"current_card_id"`
	RemainingIDs  []string   `json:"remaining_ids"`
	TimeLimitSec  int        `json:"time_limit_sec"`
	CreatedAt     time.Time  `json:"created_at"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
}

// Player represents a participant in a game.
type Player struct {
	ID           string   `json:"id"`
	DisplayName  string   `json:"display_name"`
	Score        int      `json:"score"`
	CapturedIDs  []string `json:"captured_ids"`
	IsConnected  bool     `json:"is_connected"`
	DisconnectAt *time.Time `json:"disconnect_at,omitempty"`
}
