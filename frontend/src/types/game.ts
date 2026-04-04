export type GameStatus = "waiting" | "playing" | "finished";

export interface Player {
  id: string;
  display_name: string;
  score: number;
  captured_ids: string[];
  is_connected: boolean;
}

export interface Game {
  id: string;
  room_code: string;
  status: GameStatus;
  deck_id: string;
  players: Player[];
  current_card_id: string;
  remaining_ids: string[];
  time_limit_sec: number;
}

// H4: removed "start" and "card_grabbed" that do not exist in backend ws/message.go
export type WSMessageType =
  | "join"
  | "card_revealed"
  | "grab"
  | "grab_result"
  | "game_over"
  | "player_joined"
  | "player_left"
  | "reconnect"
  | "match_found"
  | "waiting"
  | "game_state"
  | "error";

export interface WSMessage {
  type: WSMessageType;
  payload: unknown;
}

// --- Payload types matching backend ws/message.go ---

export interface JoinPayload {
  room_code: string;
  player_name: string;
  random_match: boolean;
}

export interface GrabPayload {
  card_id: string;
}

export interface ReconnectPayload {
  game_id: string;
  player_id: string;
}

export interface PlayerJoinedPayload {
  player_id: string;
  player_name: string;
  room_code: string;
}

export interface CandidateCard {
  id: string;
  name: string;
  image_url: string;
}

export interface CardRevealedPayload {
  reading_text: string;
  candidates: CandidateCard[];
  card_index: number;
  total_cards: number;
}

export interface GrabResultPayload {
  winner_id: string;
  winner_name: string;
  card_id: string;
  card_name: string;
  correct: boolean;
}

export interface PlayerResult {
  player_id: string;
  player_name: string;
  score: number;
  captured_ids: string[];
}

export interface GameOverPayload {
  players: PlayerResult[];
  winner_id: string;
}

export interface WaitingPayload {
  message: string;
}

export interface ErrorPayload {
  message: string;
}

export interface GameStatePayload {
  game_id: string;
  current_index: number;
  total_cards: number;
  status: GameStatus;
}
