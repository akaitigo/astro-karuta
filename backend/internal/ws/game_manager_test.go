package ws_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/model"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
	"github.com/akaitigo/astro-karuta/backend/internal/ws"
)

func seedTestCards(t *testing.T) *repository.InMemoryCardRepository {
	t.Helper()
	repo := repository.NewInMemoryCardRepository()
	cards := []model.Card{
		{ID: "card-1", Name: "オリオン座", Category: model.CardCategoryConstellation, ReadingText: "冬の三つ星", BestSeason: "winter"},
		{ID: "card-2", Name: "さそり座", Category: model.CardCategoryConstellation, ReadingText: "赤い心臓", BestSeason: "summer"},
		{ID: "card-3", Name: "火星", Category: model.CardCategoryPlanet, ReadingText: "赤い惑星", BestSeason: "all"},
		{ID: "card-4", Name: "木星", Category: model.CardCategoryPlanet, ReadingText: "最大の惑星", BestSeason: "all"},
		{ID: "card-5", Name: "土星", Category: model.CardCategoryPlanet, ReadingText: "環の惑星", BestSeason: "all"},
	}
	for i := range cards {
		cards[i].CreatedAt = time.Now()
		if err := repo.Create(context.Background(), &cards[i]); err != nil {
			t.Fatal(err)
		}
	}
	return repo
}

func TestGameManager_JoinAndMatchmaking(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	// Register two clients
	c1 := newTestClient("player1", "")
	c2 := newTestClient("player2", "")
	hub.Register(c1)
	hub.Register(c2)

	// Player 1 requests random match
	gm.HandleJoin("player1", ws.JoinPayload{
		PlayerName:  "Alice",
		RandomMatch: true,
	})

	// Should get a waiting message
	select {
	case msg := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgWaiting {
			t.Errorf("expected waiting, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
	}

	// Player 2 requests random match
	gm.HandleJoin("player2", ws.JoinPayload{
		PlayerName:  "Bob",
		RandomMatch: true,
	})

	// Both should get match_found
	for _, c := range []*ws.Client{c1, c2} {
		select {
		case msg := <-c.Send:
			var m ws.Message
			if err := json.Unmarshal(msg, &m); err != nil {
				t.Fatal(err)
			}
			if m.Type != ws.MsgMatchFound {
				t.Errorf("expected match_found, got %s", m.Type)
			}
		case <-time.After(time.Second):
			t.Fatalf("timeout waiting for match_found for %s", c.ID)
		}
	}

	// Game should start, so card_revealed should be sent
	for _, c := range []*ws.Client{c1, c2} {
		select {
		case msg := <-c.Send:
			var m ws.Message
			if err := json.Unmarshal(msg, &m); err != nil {
				t.Fatal(err)
			}
			if m.Type != ws.MsgCardRevealed {
				t.Errorf("expected card_revealed, got %s", m.Type)
			}
		case <-time.After(time.Second):
			t.Fatalf("timeout waiting for card_revealed for %s", c.ID)
		}
	}
}

func TestGameManager_RoomJoin(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("host", "")
	c2 := newTestClient("guest", "")
	hub.Register(c1)
	hub.Register(c2)

	// Host creates a room
	gm.HandleJoin("host", ws.JoinPayload{
		PlayerName: "Host",
		RoomCode:   "TEST01",
	})

	// Host should get player_joined
	select {
	case msg := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgPlayerJoined {
			t.Errorf("expected player_joined, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for player_joined")
	}

	// Guest joins the room
	gm.HandleJoin("guest", ws.JoinPayload{
		PlayerName: "Guest",
		RoomCode:   "TEST01",
	})

	// Both should get player_joined broadcast, then card_revealed
	for _, c := range []*ws.Client{c1, c2} {
		select {
		case msg := <-c.Send:
			var m ws.Message
			if err := json.Unmarshal(msg, &m); err != nil {
				t.Fatal(err)
			}
			if m.Type != ws.MsgPlayerJoined {
				t.Errorf("expected player_joined, got %s", m.Type)
			}
		case <-time.After(time.Second):
			t.Fatalf("timeout for %s", c.ID)
		}
	}
}

func TestGameManager_GrabCard(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	c2 := newTestClient("p2", "")
	hub.Register(c1)
	hub.Register(c2)

	// Set up a match
	gm.HandleJoin("p1", ws.JoinPayload{PlayerName: "P1", RandomMatch: true})
	<-c1.Send // waiting

	gm.HandleJoin("p2", ws.JoinPayload{PlayerName: "P2", RandomMatch: true})

	// Drain match_found
	<-c1.Send
	<-c2.Send

	// Get card_revealed to find the correct card
	msg1 := <-c1.Send
	var cardMsg ws.Message
	if err := json.Unmarshal(msg1, &cardMsg); err != nil {
		t.Fatal(err)
	}
	<-c2.Send // drain c2's card_revealed

	var revealed ws.CardRevealedPayload
	if err := json.Unmarshal(cardMsg.Payload, &revealed); err != nil {
		t.Fatal(err)
	}

	// Find the correct card (matches reading_text)
	var correctID string
	for _, c := range revealed.Candidates {
		correctID = c.ID
		break
	}

	// Player 1 grabs a card
	gm.HandleGrab("p1", ws.GrabPayload{CardID: correctID})

	// Both should get grab_result
	for _, c := range []*ws.Client{c1, c2} {
		select {
		case msg := <-c.Send:
			var m ws.Message
			if err := json.Unmarshal(msg, &m); err != nil {
				t.Fatal(err)
			}
			if m.Type != ws.MsgGrabResult {
				t.Errorf("expected grab_result, got %s", m.Type)
			}
		case <-time.After(time.Second):
			t.Fatalf("timeout for grab_result for %s", c.ID)
		}
	}
}

func TestGameManager_ProcessMessage_InvalidJSON(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	hub.Register(c1)

	gm.ProcessMessage("p1", []byte(`invalid json`))

	select {
	case msg := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgError {
			t.Errorf("expected error, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestGameManager_ProcessMessage_UnknownType(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	hub.Register(c1)

	msg, _ := ws.MarshalMessage("unknown", map[string]string{})
	gm.ProcessMessage("p1", msg)

	select {
	case raw := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgError {
			t.Errorf("expected error, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

// --- C1: PlayerName validation tests ---

func TestGameManager_Join_EmptyPlayerName(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	hub.Register(c1)

	gm.HandleJoin("p1", ws.JoinPayload{PlayerName: "", RandomMatch: true})

	select {
	case raw := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgError {
			t.Errorf("expected error for empty name, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestGameManager_Join_PlayerNameTooLong(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	hub.Register(c1)

	longName := "ABCDEFGHIJKLMNOPQRSTU" // 21 chars
	gm.HandleJoin("p1", ws.JoinPayload{PlayerName: longName, RandomMatch: true})

	select {
	case raw := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgError {
			t.Errorf("expected error for long name, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestGameManager_Join_PlayerNameControlChars(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	hub.Register(c1)

	gm.HandleJoin("p1", ws.JoinPayload{PlayerName: "bad\x00name", RandomMatch: true})

	select {
	case raw := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgError {
			t.Errorf("expected error for control chars, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

// --- C1: RoomCode validation tests ---

func TestGameManager_Join_InvalidRoomCode(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	hub.Register(c1)

	// lowercase is invalid
	gm.HandleJoin("p1", ws.JoinPayload{PlayerName: "Alice", RoomCode: "abcdef"})

	select {
	case raw := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgError {
			t.Errorf("expected error for invalid room code, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestGameManager_Join_RoomCodeWrongLength(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	hub.Register(c1)

	gm.HandleJoin("p1", ws.JoinPayload{PlayerName: "Alice", RoomCode: "ABC"})

	select {
	case raw := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgError {
			t.Errorf("expected error for short room code, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestGameManager_Join_ValidRoomCode(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	hub.Register(c1)

	gm.HandleJoin("p1", ws.JoinPayload{PlayerName: "Alice", RoomCode: "ABC123"})

	select {
	case raw := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		// Should get player_joined (not error)
		if m.Type != ws.MsgPlayerJoined {
			t.Errorf("expected player_joined, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

// --- C2: HandleGrab when not playing ---

func TestGameManager_GrabNotPlaying(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	hub.Register(c1)

	// Create a room but don't start the game (only 1 player, status=waiting)
	gm.HandleJoin("p1", ws.JoinPayload{PlayerName: "Alice", RoomCode: "TEST99"})
	<-c1.Send // drain player_joined

	// Try to grab
	gm.HandleGrab("p1", ws.GrabPayload{CardID: "card-1"})

	select {
	case raw := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgError {
			t.Errorf("expected error for grab-while-not-playing, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

// --- H2: grab limit per round ---

func TestGameManager_GrabLimitPerRound(t *testing.T) {
	hub := ws.NewHub()
	repo := seedTestCards(t)
	gm := ws.NewGameManager(hub, repo)

	c1 := newTestClient("p1", "")
	c2 := newTestClient("p2", "")
	hub.Register(c1)
	hub.Register(c2)

	// Setup match
	gm.HandleJoin("p1", ws.JoinPayload{PlayerName: "P1", RandomMatch: true})
	<-c1.Send // waiting

	gm.HandleJoin("p2", ws.JoinPayload{PlayerName: "P2", RandomMatch: true})
	<-c1.Send // match_found
	<-c2.Send // match_found

	// Get card_revealed
	msg1 := <-c1.Send
	<-c2.Send // drain c2's card_revealed

	var cardMsg ws.Message
	if err := json.Unmarshal(msg1, &cardMsg); err != nil {
		t.Fatal(err)
	}

	// First grab (wrong card on purpose)
	gm.HandleGrab("p1", ws.GrabPayload{CardID: "wrong-card-id"})

	// Drain grab_result for both
	<-c1.Send
	<-c2.Send

	// Second grab should be rejected (already grabbed this round)
	gm.HandleGrab("p1", ws.GrabPayload{CardID: "another-wrong-id"})

	select {
	case raw := <-c1.Send:
		var m ws.Message
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		if m.Type != ws.MsgError {
			t.Errorf("expected error for duplicate grab, got %s", m.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}
