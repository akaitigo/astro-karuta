export type MissionStatus = "active" | "completed" | "expired";

export interface UserMission {
  id: string;
  user_id: string;
  mission_id: string;
  card_id: string;
  title: string;
  description: string;
  status: MissionStatus;
  valid_from: string;
  valid_to: string;
  completed_at?: string;
  created_at: string;
}

export interface CompleteMissionRequest {
  user_id: string;
  lat: number;
  lng: number;
}

// H4: bonus_card matches backend model.Card (all fields from JSON serialization)
export interface CompleteMissionResponse {
  mission: UserMission;
  bonus_card?: {
    id: string;
    name: string;
    category: string;
    reading_text: string;
    image_url: string;
    description: string;
    magnitude?: number;
    distance?: string;
    best_season: string;
    created_at: string;
  };
}
