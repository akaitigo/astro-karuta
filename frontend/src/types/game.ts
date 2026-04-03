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

export type WSMessageType =
  | "join"
  | "start"
  | "card_revealed"
  | "card_grabbed"
  | "grab_result"
  | "game_over"
  | "player_joined"
  | "player_left"
  | "reconnect"
  | "error";

export interface WSMessage {
  type: WSMessageType;
  payload: unknown;
}
