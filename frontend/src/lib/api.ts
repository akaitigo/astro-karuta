import type { Card } from "@/types/card";
import type { CollectionEntry, CollectionStats } from "@/types/collection";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface ApiError {
  error: string;
}

async function fetchJSON<T>(url: string): Promise<T> {
  const res = await fetch(url);
  if (!res.ok) {
    const body: ApiError = await res.json();
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export function getCards(params?: {
  category?: string;
  season?: string;
}): Promise<Card[]> {
  const query = new URLSearchParams();
  if (params?.category) query.set("category", params.category);
  if (params?.season) query.set("season", params.season);
  const qs = query.toString();
  return fetchJSON<Card[]>(`${API_BASE}/api/v1/cards${qs ? `?${qs}` : ""}`);
}

export function getCard(id: string): Promise<Card> {
  return fetchJSON<Card>(`${API_BASE}/api/v1/cards/${encodeURIComponent(id)}`);
}

export function getCollection(params: {
  userId: string;
  category?: string;
}): Promise<CollectionEntry[]> {
  const query = new URLSearchParams();
  query.set("user_id", params.userId);
  if (params.category) query.set("category", params.category);
  return fetchJSON<CollectionEntry[]>(
    `${API_BASE}/api/v1/collections?${query.toString()}`
  );
}

export function getCollectionStats(
  userId: string
): Promise<CollectionStats> {
  const query = new URLSearchParams();
  query.set("user_id", userId);
  return fetchJSON<CollectionStats>(
    `${API_BASE}/api/v1/collections/stats?${query.toString()}`
  );
}
