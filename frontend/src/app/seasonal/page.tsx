"use client";

import { useEffect, useState, useCallback } from "react";
import type { Deck, Card } from "@/types/card";

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export default function SeasonalPage() {
  const [deck, setDeck] = useState<Deck | null>(null);
  const [cards, setCards] = useState<Card[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const currentMonth = new Date().getMonth() + 1;
  const monthNames = [
    "", "1月", "2月", "3月", "4月", "5月", "6月",
    "7月", "8月", "9月", "10月", "11月", "12月",
  ];

  const fetchSeasonalDeck = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const deckRes = await fetch(`${API_BASE}/api/v1/decks/seasonal`);
      if (!deckRes.ok) {
        throw new Error("季節デッキの取得に失敗しました");
      }
      const deckData: Deck = await deckRes.json();
      setDeck(deckData);

      // Fetch each card in the deck
      const cardPromises = deckData.card_ids.map(async (id) => {
        const res = await fetch(`${API_BASE}/api/v1/cards/${id}`);
        if (!res.ok) {
          return null;
        }
        const card: Card = await res.json();
        return card;
      });

      const fetchedCards = await Promise.all(cardPromises);
      setCards(fetchedCards.filter((c): c is Card => c !== null));
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("不明なエラーが発生しました");
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void fetchSeasonalDeck();
  }, [fetchSeasonalDeck]);

  if (loading) {
    return (
      <main style={{ padding: "24px", maxWidth: "800px", margin: "0 auto" }}>
        <h1>{monthNames[currentMonth]}の星座デッキ</h1>
        <p>読み込み中...</p>
      </main>
    );
  }

  if (error) {
    return (
      <main style={{ padding: "24px", maxWidth: "800px", margin: "0 auto" }}>
        <h1>{monthNames[currentMonth]}の星座デッキ</h1>
        <p style={{ color: "#D32F2F" }}>{error}</p>
        <button onClick={() => void fetchSeasonalDeck()}>再試行</button>
      </main>
    );
  }

  return (
    <main style={{ padding: "24px", maxWidth: "800px", margin: "0 auto" }}>
      <h1>{monthNames[currentMonth]}の星座デッキ</h1>

      {deck && (
        <p style={{ color: "#666", marginBottom: "16px" }}>
          {deck.name} - {cards.length}枚の星座カード
        </p>
      )}

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fill, minmax(240px, 1fr))",
          gap: "16px",
          marginBottom: "24px",
        }}
      >
        {cards.map((card) => (
          <div
            key={card.id}
            style={{
              border: "2px solid #1565C0",
              borderRadius: "12px",
              padding: "16px",
              backgroundColor: "#E3F2FD",
            }}
            data-testid="seasonal-card"
          >
            <h3 style={{ margin: "0 0 8px 0", color: "#1565C0" }}>
              {card.name}
            </h3>
            <p style={{ fontSize: "14px", color: "#555", margin: "0 0 8px 0" }}>
              {card.reading_text}
            </p>
            <p style={{ fontSize: "12px", color: "#888", margin: 0 }}>
              {card.description}
            </p>
            {card.magnitude !== undefined && (
              <p style={{ fontSize: "12px", color: "#999", margin: "4px 0 0 0" }}>
                等級: {card.magnitude}
              </p>
            )}
          </div>
        ))}
      </div>

      <a
        href="/"
        style={{
          display: "inline-block",
          backgroundColor: "#FF9800",
          color: "#FFFFFF",
          padding: "12px 24px",
          borderRadius: "8px",
          textDecoration: "none",
          fontWeight: "bold",
          fontSize: "16px",
        }}
      >
        このデッキで遊ぶ
      </a>
    </main>
  );
}
