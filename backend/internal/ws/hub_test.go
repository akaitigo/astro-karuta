package ws_test

import (
	"testing"

	"github.com/akaitigo/astro-karuta/backend/internal/ws"
)

func newTestClient(id, roomCode string) *ws.Client {
	return &ws.Client{
		ID:       id,
		RoomCode: roomCode,
		Send:     make(chan []byte, 256),
	}
}

func TestHub_RegisterUnregister(t *testing.T) {
	hub := ws.NewHub()
	client := newTestClient("c1", "room1")

	hub.Register(client)

	if hub.RoomSize("room1") != 1 {
		t.Errorf("expected room size 1, got %d", hub.RoomSize("room1"))
	}

	hub.Unregister(client)

	if hub.RoomSize("room1") != 0 {
		t.Errorf("expected room size 0, got %d", hub.RoomSize("room1"))
	}
}

func TestHub_BroadcastToRoom(t *testing.T) {
	hub := ws.NewHub()
	c1 := newTestClient("c1", "room1")
	c2 := newTestClient("c2", "room1")
	c3 := newTestClient("c3", "room2")

	hub.Register(c1)
	hub.Register(c2)
	hub.Register(c3)

	msg := []byte(`{"type":"test"}`)
	hub.BroadcastToRoom("room1", msg)

	// c1 and c2 should receive the message
	select {
	case got := <-c1.Send:
		if string(got) != string(msg) {
			t.Errorf("c1: expected %s, got %s", msg, got)
		}
	default:
		t.Error("c1 did not receive message")
	}

	select {
	case got := <-c2.Send:
		if string(got) != string(msg) {
			t.Errorf("c2: expected %s, got %s", msg, got)
		}
	default:
		t.Error("c2 did not receive message")
	}

	// c3 should not receive
	select {
	case <-c3.Send:
		t.Error("c3 should not receive message from room1")
	default:
		// expected
	}
}

func TestHub_SendTo(t *testing.T) {
	hub := ws.NewHub()
	c1 := newTestClient("c1", "room1")
	c2 := newTestClient("c2", "room1")

	hub.Register(c1)
	hub.Register(c2)

	msg := []byte(`{"type":"direct"}`)
	hub.SendTo("c1", msg)

	select {
	case got := <-c1.Send:
		if string(got) != string(msg) {
			t.Errorf("expected %s, got %s", msg, got)
		}
	default:
		t.Error("c1 did not receive message")
	}

	select {
	case <-c2.Send:
		t.Error("c2 should not receive direct message to c1")
	default:
		// expected
	}
}

func TestHub_JoinRoom(t *testing.T) {
	hub := ws.NewHub()
	c1 := newTestClient("c1", "room1")
	hub.Register(c1)

	if hub.RoomSize("room1") != 1 {
		t.Fatalf("expected room1 size 1")
	}

	hub.JoinRoom("c1", "room2")

	if hub.RoomSize("room1") != 0 {
		t.Errorf("expected room1 size 0 after move, got %d", hub.RoomSize("room1"))
	}
	if hub.RoomSize("room2") != 1 {
		t.Errorf("expected room2 size 1 after move, got %d", hub.RoomSize("room2"))
	}
}

func TestHub_GetClientsByRoom(t *testing.T) {
	hub := ws.NewHub()
	c1 := newTestClient("c1", "room1")
	c2 := newTestClient("c2", "room1")
	hub.Register(c1)
	hub.Register(c2)

	clients := hub.GetClientsByRoom("room1")
	if len(clients) != 2 {
		t.Errorf("expected 2 clients, got %d", len(clients))
	}

	empty := hub.GetClientsByRoom("nonexistent")
	if len(empty) != 0 {
		t.Errorf("expected 0 clients, got %d", len(empty))
	}
}
