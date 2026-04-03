package ws

import "encoding/json"

// MessageType represents the type of WebSocket message.
type MessageType string

const (
	MsgJoin         MessageType = "join"
	MsgStart        MessageType = "start"
	MsgGrab         MessageType = "grab"
	MsgReconnect    MessageType = "reconnect"
	MsgPlayerJoined MessageType = "player_joined"
	MsgPlayerLeft   MessageType = "player_left"
	MsgCardRevealed MessageType = "card_revealed"
	MsgGrabResult   MessageType = "grab_result"
	MsgGameOver     MessageType = "game_over"
	MsgError        MessageType = "error"
	MsgMatchFound   MessageType = "match_found"
	MsgWaiting      MessageType = "waiting"
	MsgGameState    MessageType = "game_state"
)

// Message is the WebSocket message envelope.
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// JoinPayload is the payload for a join message.
type JoinPayload struct {
	RoomCode    string `json:"room_code"`
	PlayerName  string `json:"player_name"`
	RandomMatch bool   `json:"random_match"`
}

// GrabPayload is sent when a player grabs a card.
type GrabPayload struct {
	CardID string `json:"card_id"`
}

// ReconnectPayload is sent to reconnect to an existing game.
type ReconnectPayload struct {
	GameID   string `json:"game_id"`
	PlayerID string `json:"player_id"`
}

// PlayerJoinedPayload notifies room that a player joined.
type PlayerJoinedPayload struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	RoomCode   string `json:"room_code"`
}

// CardRevealedPayload is sent when a new card is revealed.
type CardRevealedPayload struct {
	ReadingText string           `json:"reading_text"`
	Candidates  []CandidateCard  `json:"candidates"`
	CardIndex   int              `json:"card_index"`
	TotalCards  int              `json:"total_cards"`
}

// CandidateCard is a selectable card on the field.
type CandidateCard struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// GrabResultPayload is the result of a grab attempt.
type GrabResultPayload struct {
	WinnerID   string `json:"winner_id"`
	WinnerName string `json:"winner_name"`
	CardID     string `json:"card_id"`
	CardName   string `json:"card_name"`
	Correct    bool   `json:"correct"`
}

// GameOverPayload is sent when the game ends.
type GameOverPayload struct {
	Players  []PlayerResult `json:"players"`
	WinnerID string         `json:"winner_id"`
}

// PlayerResult contains a player's final score.
type PlayerResult struct {
	PlayerID    string `json:"player_id"`
	PlayerName  string `json:"player_name"`
	Score       int    `json:"score"`
	CapturedIDs []string `json:"captured_ids"`
}

// ErrorPayload contains an error message.
type ErrorPayload struct {
	Message string `json:"message"`
}

// MarshalMessage creates a JSON message.
func MarshalMessage(msgType MessageType, payload interface{}) ([]byte, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Message{
		Type:    msgType,
		Payload: p,
	})
}
