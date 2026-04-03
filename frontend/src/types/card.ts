export type CardCategory = "constellation" | "planet" | "phenomenon";

export interface Card {
  id: string;
  name: string;
  category: CardCategory;
  reading_text: string;
  image_url: string;
  description: string;
  magnitude?: number;
  distance?: string;
  best_season: string;
}

export interface Deck {
  id: string;
  name: string;
  card_ids: string[];
  seasonal: boolean;
  valid_from: string;
  valid_to: string;
}
