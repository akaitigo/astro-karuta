export interface CollectionEntry {
  id: string;
  user_id: string;
  card_id: string;
  obtained_at: string;
  source: "game" | "mission";
}

export interface CollectionStats {
  user_id: string;
  total_cards: number;
  collected: number;
  percentage: number;
}

export interface ObservationMission {
  id: string;
  card_id: string;
  title: string;
  valid_from: string;
  valid_to: string;
}
