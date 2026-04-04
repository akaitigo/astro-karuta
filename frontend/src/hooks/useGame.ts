"use client";

import { useCallback, useMemo, useReducer } from "react";
import { useWebSocket, type ConnectionStatus } from "./useWebSocket";
import type {
  WSMessage,
  GameStatus,
  CandidateCard,
  CardRevealedPayload,
  GrabResultPayload,
  GameOverPayload,
  GameStatePayload,
  PlayerJoinedPayload,
  PlayerResult,
  WaitingPayload,
  ErrorPayload,
} from "@/types/game";

// --- State ---

export interface GameState {
  gameStatus: GameStatus | "lobby" | "matchmaking";
  roomCode: string;
  playerId: string;
  playerName: string;
  readingText: string;
  candidates: CandidateCard[];
  cardIndex: number;
  totalCards: number;
  scores: Record<string, number>;
  playerNames: Record<string, string>;
  lastGrabResult: GrabResultPayload | null;
  gameResult: GameOverPayload | null;
  error: string | null;
}

const initialState: GameState = {
  gameStatus: "lobby",
  roomCode: "",
  playerId: "",
  playerName: "",
  readingText: "",
  candidates: [],
  cardIndex: 0,
  totalCards: 0,
  scores: {},
  playerNames: {},
  lastGrabResult: null,
  gameResult: null,
  error: null,
};

// --- Actions ---

type GameAction =
  | { type: "SET_PLAYER"; playerId: string; playerName: string }
  | { type: "START_MATCHMAKING" }
  | { type: "PLAYER_JOINED"; payload: PlayerJoinedPayload }
  | { type: "CARD_REVEALED"; payload: CardRevealedPayload }
  | { type: "GRAB_RESULT"; payload: GrabResultPayload }
  | { type: "GAME_OVER"; payload: GameOverPayload }
  | { type: "GAME_STATE"; payload: GameStatePayload }
  | { type: "WAITING"; payload: WaitingPayload }
  | { type: "MATCH_FOUND"; roomCode: string }
  | { type: "ERROR"; message: string }
  | { type: "RESET" };

function gameReducer(state: GameState, action: GameAction): GameState {
  switch (action.type) {
    case "SET_PLAYER":
      return {
        ...state,
        playerId: action.playerId,
        playerName: action.playerName,
      };

    case "START_MATCHMAKING":
      return { ...state, gameStatus: "matchmaking", error: null };

    case "PLAYER_JOINED": {
      const { player_id, player_name, room_code } = action.payload;
      return {
        ...state,
        roomCode: room_code,
        playerNames: { ...state.playerNames, [player_id]: player_name },
        scores: { ...state.scores, [player_id]: 0 },
      };
    }

    case "CARD_REVEALED":
      return {
        ...state,
        gameStatus: "playing",
        readingText: action.payload.reading_text,
        candidates: action.payload.candidates,
        cardIndex: action.payload.card_index,
        totalCards: action.payload.total_cards,
        lastGrabResult: null,
      };

    case "GRAB_RESULT": {
      const result = action.payload;
      const newScores = { ...state.scores };
      if (result.correct && result.winner_id) {
        newScores[result.winner_id] = (newScores[result.winner_id] ?? 0) + 1;
      }
      // Update player name mapping if we learn a new name
      const newNames = { ...state.playerNames };
      if (result.winner_name && result.winner_id) {
        newNames[result.winner_id] = result.winner_name;
      }
      return {
        ...state,
        scores: newScores,
        playerNames: newNames,
        lastGrabResult: result,
      };
    }

    case "GAME_OVER": {
      const nameMap: Record<string, string> = { ...state.playerNames };
      const scoreMap: Record<string, number> = { ...state.scores };
      for (const p of action.payload.players) {
        nameMap[p.player_id] = p.player_name;
        scoreMap[p.player_id] = p.score;
      }
      return {
        ...state,
        gameStatus: "finished",
        gameResult: action.payload,
        playerNames: nameMap,
        scores: scoreMap,
      };
    }

    // H4: handle game_state message from reconnect
    case "GAME_STATE":
      return {
        ...state,
        gameStatus: action.payload.status,
        cardIndex: action.payload.current_index,
        totalCards: action.payload.total_cards,
      };

    case "WAITING":
      return { ...state, gameStatus: "matchmaking" };

    case "MATCH_FOUND":
      return { ...state, roomCode: action.roomCode, gameStatus: "waiting" };

    case "ERROR":
      return { ...state, error: action.message };

    case "RESET":
      return initialState;

    default:
      return state;
  }
}

// --- Hook ---

interface UseGameOptions {
  wsUrl: string;
}

interface UseGameReturn {
  state: GameState;
  connectionStatus: ConnectionStatus;
  joinRoom: (roomCode: string, playerName: string) => void;
  createRoom: (playerName: string) => void;
  randomMatch: (playerName: string) => void;
  grabCard: (cardId: string) => void;
  reset: () => void;
  players: PlayerResult[];
}

function isCardRevealedPayload(v: unknown): v is CardRevealedPayload {
  if (typeof v !== "object" || v === null) return false;
  const obj = v as Record<string, unknown>;
  return (
    typeof obj.reading_text === "string" &&
    Array.isArray(obj.candidates) &&
    typeof obj.card_index === "number" &&
    typeof obj.total_cards === "number"
  );
}

function isGrabResultPayload(v: unknown): v is GrabResultPayload {
  if (typeof v !== "object" || v === null) return false;
  const obj = v as Record<string, unknown>;
  return typeof obj.card_id === "string" && typeof obj.correct === "boolean";
}

function isGameOverPayload(v: unknown): v is GameOverPayload {
  if (typeof v !== "object" || v === null) return false;
  const obj = v as Record<string, unknown>;
  return Array.isArray(obj.players) && typeof obj.winner_id === "string";
}

function isPlayerJoinedPayload(v: unknown): v is PlayerJoinedPayload {
  if (typeof v !== "object" || v === null) return false;
  const obj = v as Record<string, unknown>;
  return (
    typeof obj.player_id === "string" && typeof obj.room_code === "string"
  );
}

function isWaitingPayload(v: unknown): v is WaitingPayload {
  if (typeof v !== "object" || v === null) return false;
  const obj = v as Record<string, unknown>;
  return typeof obj.message === "string";
}

function isErrorPayload(v: unknown): v is ErrorPayload {
  if (typeof v !== "object" || v === null) return false;
  const obj = v as Record<string, unknown>;
  return typeof obj.message === "string";
}

function isGameStatePayload(v: unknown): v is GameStatePayload {
  if (typeof v !== "object" || v === null) return false;
  const obj = v as Record<string, unknown>;
  return (
    typeof obj.game_id === "string" &&
    typeof obj.current_index === "number" &&
    typeof obj.total_cards === "number" &&
    typeof obj.status === "string"
  );
}

function isMatchFoundPayload(
  v: unknown,
): v is { room_code: string } {
  if (typeof v !== "object" || v === null) return false;
  const obj = v as Record<string, unknown>;
  return typeof obj.room_code === "string";
}

export function useGame({ wsUrl }: UseGameOptions): UseGameReturn {
  const [state, dispatch] = useReducer(gameReducer, initialState);

  const handleMessage = useCallback((message: WSMessage) => {
    switch (message.type) {
      case "player_joined":
        if (isPlayerJoinedPayload(message.payload)) {
          dispatch({ type: "PLAYER_JOINED", payload: message.payload });
        }
        break;

      case "card_revealed":
        if (isCardRevealedPayload(message.payload)) {
          dispatch({ type: "CARD_REVEALED", payload: message.payload });
        }
        break;

      case "grab_result":
        if (isGrabResultPayload(message.payload)) {
          dispatch({ type: "GRAB_RESULT", payload: message.payload });
        }
        break;

      case "game_over":
        if (isGameOverPayload(message.payload)) {
          dispatch({ type: "GAME_OVER", payload: message.payload });
        }
        break;

      // H4: handle game_state from reconnect
      case "game_state":
        if (isGameStatePayload(message.payload)) {
          dispatch({ type: "GAME_STATE", payload: message.payload });
        }
        break;

      case "waiting":
        if (isWaitingPayload(message.payload)) {
          dispatch({ type: "WAITING", payload: message.payload });
        }
        break;

      case "match_found":
        if (isMatchFoundPayload(message.payload)) {
          dispatch({
            type: "MATCH_FOUND",
            roomCode: message.payload.room_code,
          });
        }
        break;

      case "error":
        if (isErrorPayload(message.payload)) {
          dispatch({ type: "ERROR", message: message.payload.message });
        }
        break;
    }
  }, []);

  const { status, send, connect, disconnect, waitForConnection } = useWebSocket({
    url: wsUrl,
    onMessage: handleMessage,
    autoConnect: false,
  });

  // H9: helper to connect and wait for open before sending
  const connectAndSend = useCallback(
    (payload: { room_code: string; player_name: string; random_match: boolean }) => {
      connect();
      waitForConnection()
        .then(() => {
          send("join", payload);
        })
        .catch(() => {
          // Connection was closed before open (e.g. disconnect/reset called);
          // no action needed since disconnect already handles cleanup.
        });
    },
    [connect, waitForConnection, send],
  );

  const joinRoom = useCallback(
    (roomCode: string, playerName: string) => {
      dispatch({
        type: "SET_PLAYER",
        playerId: "",
        playerName,
      });
      connectAndSend({
        room_code: roomCode,
        player_name: playerName,
        random_match: false,
      });
    },
    [connectAndSend],
  );

  const createRoom = useCallback(
    (playerName: string) => {
      dispatch({
        type: "SET_PLAYER",
        playerId: "",
        playerName,
      });
      connectAndSend({
        room_code: "",
        player_name: playerName,
        random_match: false,
      });
    },
    [connectAndSend],
  );

  const randomMatch = useCallback(
    (playerName: string) => {
      dispatch({
        type: "SET_PLAYER",
        playerId: "",
        playerName,
      });
      dispatch({ type: "START_MATCHMAKING" });
      connectAndSend({
        room_code: "",
        player_name: playerName,
        random_match: true,
      });
    },
    [connectAndSend],
  );

  const grabCard = useCallback(
    (cardId: string) => {
      send("grab", { card_id: cardId });
    },
    [send],
  );

  const reset = useCallback(() => {
    disconnect();
    dispatch({ type: "RESET" });
  }, [disconnect]);

  const players: PlayerResult[] = useMemo(() => {
    return Object.entries(state.scores).map(([id, score]) => ({
      player_id: id,
      player_name: state.playerNames[id] ?? "Unknown",
      score,
      captured_ids: [],
    }));
  }, [state.scores, state.playerNames]);

  return {
    state,
    connectionStatus: status,
    joinRoom,
    createRoom,
    randomMatch,
    grabCard,
    reset,
    players,
  };
}
