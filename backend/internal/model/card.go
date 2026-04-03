package model

import "time"

// CardCategory represents the type of astronomical object.
type CardCategory string

const (
	CardCategoryConstellation CardCategory = "constellation"
	CardCategoryPlanet        CardCategory = "planet"
	CardCategoryPhenomenon    CardCategory = "phenomenon"
)

// Card represents a karuta card with astronomical data.
type Card struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Category    CardCategory `json:"category"`
	ReadingText string       `json:"reading_text"`
	ImageURL    string       `json:"image_url"`
	Description string       `json:"description"`
	Magnitude   *float64     `json:"magnitude,omitempty"`
	Distance    *string      `json:"distance,omitempty"`
	BestSeason  string       `json:"best_season"`
	CreatedAt   time.Time    `json:"created_at"`
}

// Deck represents a collection of cards for a game session.
type Deck struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CardIDs   []string  `json:"card_ids"`
	Seasonal  bool      `json:"seasonal"`
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
}
