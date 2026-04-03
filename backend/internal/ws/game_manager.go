package ws

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
	"github.com/google/uuid"
)

const (
	maxPlayersPerRoom = 2
	reconnectTimeout  = 30 * time.Second
	defaultTimeLimit  = 300 // 5 minutes
)

// GameState holds the state of an active game.
type GameState struct {
	mu             sync.RWMutex
	ID             string
	RoomCode       string
	Status         model.GameStatus
	Players        map[string]*PlayerState
	Cards          []model.Card
	CurrentIndex   int
	CurrentCard    *model.Card
	Candidates     []model.Card
	GrabLock       sync.Mutex
	GrabHandled    bool
	TimeLimitSec   int
	StartedAt      time.Time
}

// PlayerState tracks an individual player within a game.
type PlayerState struct {
	ID           string
	Name         string
	ClientID     string
	Score        int
	CapturedIDs  []string
	IsConnected  bool
	DisconnectAt *time.Time
}

// GameManager coordinates matchmaking and game lifecycle.
type GameManager struct {
	mu           sync.RWMutex
	hub          *Hub
	cardRepo     repository.CardRepository
	games        map[string]*GameState // gameID -> state
	roomToGame   map[string]string     // roomCode -> gameID
	waitingQueue []string              // clientIDs waiting for random match
	waitingNames map[string]string     // clientID -> playerName
}

// NewGameManager creates a new GameManager.
func NewGameManager(hub *Hub, cardRepo repository.CardRepository) *GameManager {
	return &GameManager{
		hub:          hub,
		cardRepo:     cardRepo,
		games:        make(map[string]*GameState),
		roomToGame:   make(map[string]string),
		waitingNames: make(map[string]string),
	}
}

// HandleJoin processes a player's join request.
func (gm *GameManager) HandleJoin(clientID string, payload JoinPayload) {
	if payload.RandomMatch {
		gm.handleRandomMatch(clientID, payload.PlayerName)
		return
	}

	roomCode := payload.RoomCode
	if roomCode == "" {
		roomCode = generateRoomCode()
	}

	gm.mu.Lock()
	gameID, exists := gm.roomToGame[roomCode]
	gm.mu.Unlock()

	if exists {
		gm.joinExistingGame(clientID, gameID, payload.PlayerName)
		return
	}

	gm.createGame(clientID, roomCode, payload.PlayerName)
}

func (gm *GameManager) handleRandomMatch(clientID, playerName string) {
	gm.mu.Lock()

	if len(gm.waitingQueue) > 0 {
		otherClientID := gm.waitingQueue[0]
		gm.waitingQueue = gm.waitingQueue[1:]
		otherName := gm.waitingNames[otherClientID]
		delete(gm.waitingNames, otherClientID)
		gm.mu.Unlock()

		roomCode := generateRoomCode()
		game := gm.createGameState(roomCode)

		gm.addPlayerToGame(game, otherClientID, otherName)
		gm.addPlayerToGame(game, clientID, playerName)

		gm.mu.Lock()
		gm.games[game.ID] = game
		gm.roomToGame[roomCode] = game.ID
		gm.mu.Unlock()

		gm.hub.JoinRoom(otherClientID, roomCode)
		gm.hub.JoinRoom(clientID, roomCode)

		matchMsg, _ := MarshalMessage(MsgMatchFound, PlayerJoinedPayload{
			RoomCode: roomCode,
		})
		gm.hub.SendTo(otherClientID, matchMsg)
		gm.hub.SendTo(clientID, matchMsg)

		gm.startGame(game)
		return
	}

	gm.waitingQueue = append(gm.waitingQueue, clientID)
	gm.waitingNames[clientID] = playerName
	gm.mu.Unlock()

	waitMsg, _ := MarshalMessage(MsgWaiting, map[string]string{
		"message": "Waiting for opponent...",
	})
	gm.hub.SendTo(clientID, waitMsg)
}

func (gm *GameManager) createGame(clientID, roomCode, playerName string) {
	game := gm.createGameState(roomCode)
	gm.addPlayerToGame(game, clientID, playerName)

	gm.mu.Lock()
	gm.games[game.ID] = game
	gm.roomToGame[roomCode] = game.ID
	gm.mu.Unlock()

	gm.hub.JoinRoom(clientID, roomCode)

	joinedMsg, _ := MarshalMessage(MsgPlayerJoined, PlayerJoinedPayload{
		PlayerID:   clientID,
		PlayerName: playerName,
		RoomCode:   roomCode,
	})
	gm.hub.SendTo(clientID, joinedMsg)
}

func (gm *GameManager) joinExistingGame(clientID, gameID, playerName string) {
	gm.mu.RLock()
	game, ok := gm.games[gameID]
	gm.mu.RUnlock()
	if !ok {
		gm.sendError(clientID, "game not found")
		return
	}

	game.mu.Lock()
	playerCount := len(game.Players)
	game.mu.Unlock()

	if playerCount >= maxPlayersPerRoom {
		gm.sendError(clientID, "room is full")
		return
	}

	gm.addPlayerToGame(game, clientID, playerName)
	gm.hub.JoinRoom(clientID, game.RoomCode)

	joinedMsg, _ := MarshalMessage(MsgPlayerJoined, PlayerJoinedPayload{
		PlayerID:   clientID,
		PlayerName: playerName,
		RoomCode:   game.RoomCode,
	})
	gm.hub.BroadcastToRoom(game.RoomCode, joinedMsg)

	game.mu.RLock()
	pc := len(game.Players)
	game.mu.RUnlock()

	if pc == maxPlayersPerRoom {
		gm.startGame(game)
	}
}

func (gm *GameManager) createGameState(roomCode string) *GameState {
	cards, err := gm.cardRepo.List(context.Background(), repository.CardFilter{})
	if err != nil {
		log.Printf("failed to load cards: %v", err)
		return nil
	}

	shuffled := make([]model.Card, len(cards))
	copy(shuffled, cards)
	mrand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Use 20 cards for a game
	gameCards := shuffled
	if len(gameCards) > 20 {
		gameCards = gameCards[:20]
	}

	return &GameState{
		ID:           uuid.New().String(),
		RoomCode:     roomCode,
		Status:       model.GameStatusWaiting,
		Players:      make(map[string]*PlayerState),
		Cards:        gameCards,
		CurrentIndex: 0,
		TimeLimitSec: defaultTimeLimit,
	}
}

func (gm *GameManager) addPlayerToGame(game *GameState, clientID, name string) {
	game.mu.Lock()
	defer game.mu.Unlock()

	game.Players[clientID] = &PlayerState{
		ID:          clientID,
		Name:        name,
		ClientID:    clientID,
		Score:       0,
		CapturedIDs: []string{},
		IsConnected: true,
	}
}

func (gm *GameManager) startGame(game *GameState) {
	game.mu.Lock()
	game.Status = model.GameStatusPlaying
	game.StartedAt = time.Now()
	game.mu.Unlock()

	log.Printf("game %s started in room %s", game.ID, game.RoomCode)
	gm.revealNextCard(game)
}

func (gm *GameManager) revealNextCard(game *GameState) {
	game.mu.Lock()
	if game.CurrentIndex >= len(game.Cards) {
		game.mu.Unlock()
		gm.endGame(game)
		return
	}

	card := game.Cards[game.CurrentIndex]
	game.CurrentCard = &card
	game.GrabHandled = false

	// Build candidate list: correct card + 3 distractors
	candidates := []model.Card{card}
	for i, c := range game.Cards {
		if i != game.CurrentIndex && len(candidates) < 4 {
			candidates = append(candidates, c)
		}
	}
	mrand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	game.Candidates = candidates

	cardIndex := game.CurrentIndex
	totalCards := len(game.Cards)
	game.mu.Unlock()

	candidateCards := make([]CandidateCard, len(candidates))
	for i, c := range candidates {
		candidateCards[i] = CandidateCard{
			ID:       c.ID,
			Name:     c.Name,
			ImageURL: c.ImageURL,
		}
	}

	msg, _ := MarshalMessage(MsgCardRevealed, CardRevealedPayload{
		ReadingText: card.ReadingText,
		Candidates:  candidateCards,
		CardIndex:   cardIndex + 1,
		TotalCards:  totalCards,
	})
	gm.hub.BroadcastToRoom(game.RoomCode, msg)
}

// HandleGrab processes a card grab attempt.
func (gm *GameManager) HandleGrab(clientID string, payload GrabPayload) {
	game := gm.findGameByClient(clientID)
	if game == nil {
		gm.sendError(clientID, "not in a game")
		return
	}

	game.GrabLock.Lock()
	defer game.GrabLock.Unlock()

	game.mu.RLock()
	if game.GrabHandled {
		game.mu.RUnlock()
		return
	}
	currentCard := game.CurrentCard
	game.mu.RUnlock()

	if currentCard == nil {
		return
	}

	correct := payload.CardID == currentCard.ID

	game.mu.Lock()
	game.GrabHandled = true

	playerName := ""
	if correct {
		if player, ok := game.Players[clientID]; ok {
			player.Score++
			player.CapturedIDs = append(player.CapturedIDs, currentCard.ID)
			playerName = player.Name
		}
	}
	game.CurrentIndex++
	game.mu.Unlock()

	resultMsg, _ := MarshalMessage(MsgGrabResult, GrabResultPayload{
		WinnerID:   clientID,
		WinnerName: playerName,
		CardID:     currentCard.ID,
		CardName:   currentCard.Name,
		Correct:    correct,
	})
	gm.hub.BroadcastToRoom(game.RoomCode, resultMsg)

	if correct {
		// Brief pause then reveal next card
		time.AfterFunc(2*time.Second, func() {
			gm.revealNextCard(game)
		})
	}
}

// HandleReconnect processes a reconnect attempt.
func (gm *GameManager) HandleReconnect(clientID string, payload ReconnectPayload) {
	gm.mu.RLock()
	game, ok := gm.games[payload.GameID]
	gm.mu.RUnlock()

	if !ok {
		gm.sendError(clientID, "game not found")
		return
	}

	game.mu.Lock()
	player, exists := game.Players[payload.PlayerID]
	if !exists {
		game.mu.Unlock()
		gm.sendError(clientID, "player not found in game")
		return
	}

	player.IsConnected = true
	player.DisconnectAt = nil
	player.ClientID = clientID
	game.mu.Unlock()

	gm.hub.JoinRoom(clientID, game.RoomCode)

	// Send current game state
	game.mu.RLock()
	statePayload := map[string]interface{}{
		"game_id":       game.ID,
		"current_index": game.CurrentIndex,
		"total_cards":   len(game.Cards),
		"status":        game.Status,
	}
	game.mu.RUnlock()

	stateMsg, _ := MarshalMessage(MsgGameState, statePayload)
	gm.hub.SendTo(clientID, stateMsg)
}

// HandleDisconnect handles a player disconnection.
func (gm *GameManager) HandleDisconnect(clientID string) {
	game := gm.findGameByClient(clientID)
	if game == nil {
		return
	}

	game.mu.Lock()
	for _, player := range game.Players {
		if player.ClientID == clientID {
			player.IsConnected = false
			now := time.Now()
			player.DisconnectAt = &now
			break
		}
	}
	game.mu.Unlock()

	// Notify other players
	msg, _ := MarshalMessage(MsgPlayerLeft, map[string]string{
		"player_id": clientID,
	})
	gm.hub.BroadcastToRoom(game.RoomCode, msg)

	// Set reconnect timeout
	time.AfterFunc(reconnectTimeout, func() {
		game.mu.RLock()
		allDisconnected := true
		for _, player := range game.Players {
			if player.IsConnected {
				allDisconnected = false
				break
			}
		}
		game.mu.RUnlock()

		if allDisconnected {
			gm.endGame(game)
		}
	})
}

func (gm *GameManager) endGame(game *GameState) {
	game.mu.Lock()
	if game.Status == model.GameStatusFinished {
		game.mu.Unlock()
		return
	}
	game.Status = model.GameStatusFinished

	var results []PlayerResult
	var winnerID string
	maxScore := -1
	for _, p := range game.Players {
		results = append(results, PlayerResult{
			PlayerID:    p.ID,
			PlayerName:  p.Name,
			Score:       p.Score,
			CapturedIDs: p.CapturedIDs,
		})
		if p.Score > maxScore {
			maxScore = p.Score
			winnerID = p.ID
		}
	}
	roomCode := game.RoomCode
	game.mu.Unlock()

	msg, _ := MarshalMessage(MsgGameOver, GameOverPayload{
		Players:  results,
		WinnerID: winnerID,
	})
	gm.hub.BroadcastToRoom(roomCode, msg)

	// Cleanup
	gm.mu.Lock()
	delete(gm.roomToGame, roomCode)
	delete(gm.games, game.ID)
	gm.mu.Unlock()
}

func (gm *GameManager) findGameByClient(clientID string) *GameState {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	for _, game := range gm.games {
		game.mu.RLock()
		_, ok := game.Players[clientID]
		game.mu.RUnlock()
		if ok {
			return game
		}
	}
	return nil
}

func (gm *GameManager) sendError(clientID, message string) {
	msg, _ := MarshalMessage(MsgError, ErrorPayload{Message: message})
	gm.hub.SendTo(clientID, msg)
}

func generateRoomCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			code[i] = chars[0]
			continue
		}
		code[i] = chars[n.Int64()]
	}
	return string(code)
}

// ProcessMessage routes an incoming message to the appropriate handler.
func (gm *GameManager) ProcessMessage(clientID string, data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		gm.sendError(clientID, fmt.Sprintf("invalid message: %v", err))
		return
	}

	switch msg.Type {
	case MsgJoin:
		var payload JoinPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			gm.sendError(clientID, "invalid join payload")
			return
		}
		gm.HandleJoin(clientID, payload)

	case MsgGrab:
		var payload GrabPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			gm.sendError(clientID, "invalid grab payload")
			return
		}
		gm.HandleGrab(clientID, payload)

	case MsgReconnect:
		var payload ReconnectPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			gm.sendError(clientID, "invalid reconnect payload")
			return
		}
		gm.HandleReconnect(clientID, payload)

	default:
		gm.sendError(clientID, fmt.Sprintf("unknown message type: %s", msg.Type))
	}
}
