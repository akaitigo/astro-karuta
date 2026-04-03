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

export interface CompleteMissionResponse {
  mission: UserMission;
  bonus_card?: {
    id: string;
    name: string;
    category: string;
    reading_text: string;
    description: string;
  };
}
