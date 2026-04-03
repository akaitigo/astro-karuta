import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useGame } from "./useGame";

// --- Mock WebSocket ---
type WSCallback = (event: { data: string }) => void;

class MockWebSocket {
  static instances: MockWebSocket[] = [];
  readyState: number = 0; // CONNECTING
  onopen: (() => void) | null = null;
  onclose: (() => void) | null = null;
  onmessage: WSCallback | null = null;
  onerror: (() => void) | null = null;
  sentMessages: string[] = [];

  constructor(_url: string) {
    MockWebSocket.instances.push(this);
    // Simulate async connection
    setTimeout(() => {
      this.readyState = 1; // OPEN
      this.onopen?.();
    }, 0);
  }

  send(data: string) {
    this.sentMessages.push(data);
  }

  close() {
    this.readyState = 3; // CLOSED
    this.onclose?.();
  }

  simulateMessage(type: string, payload: unknown) {
    this.onmessage?.({ data: JSON.stringify({ type, payload }) });
  }

  static get OPEN() {
    return 1;
  }
  static get CLOSED() {
    return 3;
  }
  static get CONNECTING() {
    return 0;
  }
  static get CLOSING() {
    return 2;
  }
}

describe("useGame", () => {
  beforeEach(() => {
    MockWebSocket.instances = [];
    vi.stubGlobal("WebSocket", MockWebSocket);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("initializes with lobby status", () => {
    const { result } = renderHook(() =>
      useGame({ wsUrl: "ws://test/ws" }),
    );
    expect(result.current.state.gameStatus).toBe("lobby");
  });

  it("transitions to matchmaking on randomMatch", () => {
    const { result } = renderHook(() =>
      useGame({ wsUrl: "ws://test/ws" }),
    );

    act(() => {
      result.current.randomMatch("TestPlayer");
    });

    expect(result.current.state.gameStatus).toBe("matchmaking");
    expect(result.current.state.playerName).toBe("TestPlayer");
  });

  it("handles card_revealed message", async () => {
    const { result } = renderHook(() =>
      useGame({ wsUrl: "ws://test/ws" }),
    );

    // Connect
    act(() => {
      result.current.randomMatch("TestPlayer");
    });

    // Wait for connection
    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBeGreaterThan(0);
    });

    const ws = MockWebSocket.instances[0];

    act(() => {
      ws.simulateMessage("card_revealed", {
        reading_text: "冬の大三角を形成する星座",
        candidates: [
          { id: "c1", name: "オリオン座", image_url: "/orion.png" },
          { id: "c2", name: "こいぬ座", image_url: "/canis-minor.png" },
        ],
        card_index: 1,
        total_cards: 20,
      });
    });

    expect(result.current.state.gameStatus).toBe("playing");
    expect(result.current.state.readingText).toBe(
      "冬の大三角を形成する星座",
    );
    expect(result.current.state.candidates).toHaveLength(2);
    expect(result.current.state.cardIndex).toBe(1);
    expect(result.current.state.totalCards).toBe(20);
  });

  it("handles grab_result message", async () => {
    const { result } = renderHook(() =>
      useGame({ wsUrl: "ws://test/ws" }),
    );

    act(() => {
      result.current.randomMatch("TestPlayer");
    });

    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBeGreaterThan(0);
    });

    const ws = MockWebSocket.instances[0];

    // First reveal a card
    act(() => {
      ws.simulateMessage("card_revealed", {
        reading_text: "test",
        candidates: [{ id: "c1", name: "Test", image_url: "/test.png" }],
        card_index: 1,
        total_cards: 10,
      });
    });

    // Then simulate a grab result
    act(() => {
      ws.simulateMessage("grab_result", {
        winner_id: "player-1",
        winner_name: "TestPlayer",
        card_id: "c1",
        card_name: "Test",
        correct: true,
      });
    });

    expect(result.current.state.lastGrabResult).not.toBeNull();
    expect(result.current.state.lastGrabResult?.correct).toBe(true);
    expect(result.current.state.scores["player-1"]).toBe(1);
  });

  it("handles game_over message", async () => {
    const { result } = renderHook(() =>
      useGame({ wsUrl: "ws://test/ws" }),
    );

    act(() => {
      result.current.randomMatch("TestPlayer");
    });

    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBeGreaterThan(0);
    });

    const ws = MockWebSocket.instances[0];

    act(() => {
      ws.simulateMessage("game_over", {
        players: [
          {
            player_id: "p1",
            player_name: "Alice",
            score: 8,
            captured_ids: ["c1"],
          },
          {
            player_id: "p2",
            player_name: "Bob",
            score: 5,
            captured_ids: ["c2"],
          },
        ],
        winner_id: "p1",
      });
    });

    expect(result.current.state.gameStatus).toBe("finished");
    expect(result.current.state.gameResult).not.toBeNull();
    expect(result.current.state.gameResult?.winner_id).toBe("p1");
  });

  it("resets state on reset()", async () => {
    const { result } = renderHook(() =>
      useGame({ wsUrl: "ws://test/ws" }),
    );

    act(() => {
      result.current.randomMatch("TestPlayer");
    });

    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBeGreaterThan(0);
    });

    const ws = MockWebSocket.instances[0];

    act(() => {
      ws.simulateMessage("card_revealed", {
        reading_text: "test",
        candidates: [],
        card_index: 1,
        total_cards: 10,
      });
    });

    expect(result.current.state.gameStatus).toBe("playing");

    act(() => {
      result.current.reset();
    });

    expect(result.current.state.gameStatus).toBe("lobby");
    expect(result.current.state.readingText).toBe("");
  });

  it("handles player_joined message", async () => {
    const { result } = renderHook(() =>
      useGame({ wsUrl: "ws://test/ws" }),
    );

    act(() => {
      result.current.createRoom("Host");
    });

    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBeGreaterThan(0);
    });

    const ws = MockWebSocket.instances[0];

    act(() => {
      ws.simulateMessage("player_joined", {
        player_id: "p1",
        player_name: "Alice",
        room_code: "ABC123",
      });
    });

    expect(result.current.state.roomCode).toBe("ABC123");
    expect(result.current.state.playerNames["p1"]).toBe("Alice");
  });

  it("handles error message", async () => {
    const { result } = renderHook(() =>
      useGame({ wsUrl: "ws://test/ws" }),
    );

    act(() => {
      result.current.randomMatch("TestPlayer");
    });

    await vi.waitFor(() => {
      expect(MockWebSocket.instances.length).toBeGreaterThan(0);
    });

    const ws = MockWebSocket.instances[0];

    act(() => {
      ws.simulateMessage("error", {
        message: "Room is full",
      });
    });

    expect(result.current.state.error).toBe("Room is full");
  });
});
