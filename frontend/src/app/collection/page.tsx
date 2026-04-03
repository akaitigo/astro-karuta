"use client";

import { useCallback, useEffect, useState } from "react";
import type { Card, CardCategory } from "@/types/card";
import type { CollectionEntry, CollectionStats } from "@/types/collection";
import { getCards, getCollection, getCollectionStats } from "@/lib/api";
import CollectionCard from "@/components/CollectionCard";
import CardDetail from "@/components/CardDetail";

type CategoryFilter = CardCategory | "all";

const categories: { value: CategoryFilter; label: string }[] = [
  { value: "all", label: "すべて" },
  { value: "constellation", label: "星座" },
  { value: "planet", label: "惑星" },
  { value: "phenomenon", label: "天文現象" },
];

// Temporary user ID until auth is implemented
const DEFAULT_USER_ID = "user-1";

export default function CollectionPage() {
  const [cards, setCards] = useState<Card[]>([]);
  const [collectedCardIds, setCollectedCardIds] = useState<Set<string>>(
    new Set()
  );
  const [stats, setStats] = useState<CollectionStats | null>(null);
  const [filter, setFilter] = useState<CategoryFilter>("all");
  const [selectedCard, setSelectedCard] = useState<Card | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const categoryParam = filter === "all" ? undefined : filter;

      const [cardsData, collectionData, statsData] = await Promise.all([
        getCards({ category: categoryParam }),
        getCollection({
          userId: DEFAULT_USER_ID,
          category: categoryParam,
        }),
        getCollectionStats(DEFAULT_USER_ID),
      ]);

      setCards(cardsData);
      setCollectedCardIds(
        new Set(collectionData.map((e: CollectionEntry) => e.card_id))
      );
      setStats(statsData);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "データの取得に失敗しました";
      setError(message);
    } finally {
      setLoading(false);
    }
  }, [filter]);

  useEffect(() => {
    void loadData();
  }, [loadData]);

  const handleSelectCard = (card: Card) => {
    setSelectedCard(card);
  };

  const handleCloseDetail = () => {
    setSelectedCard(null);
  };

  return (
    <main style={{ maxWidth: "960px", margin: "0 auto", padding: "24px" }}>
      <h1 style={{ fontSize: "24px", marginBottom: "8px" }}>天体コレクション</h1>

      {stats && (
        <div style={{ marginBottom: "24px" }}>
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              marginBottom: "4px",
              fontSize: "14px",
            }}
          >
            <span>
              {stats.collected} / {stats.total_cards} 収集済み
            </span>
            <span>{stats.percentage.toFixed(1)}%</span>
          </div>
          <div
            role="progressbar"
            aria-valuenow={stats.percentage}
            aria-valuemin={0}
            aria-valuemax={100}
            aria-label="コレクション進捗"
            style={{
              width: "100%",
              height: "12px",
              background: "#e0e0e0",
              borderRadius: "6px",
              overflow: "hidden",
            }}
          >
            <div
              style={{
                width: `${stats.percentage}%`,
                height: "100%",
                background: "linear-gradient(90deg, #1a237e, #4a148c)",
                borderRadius: "6px",
                transition: "width 0.3s ease",
              }}
            />
          </div>
        </div>
      )}

      <nav
        aria-label="カテゴリフィルター"
        style={{
          display: "flex",
          gap: "8px",
          marginBottom: "24px",
          flexWrap: "wrap",
        }}
      >
        {categories.map((cat) => (
          <button
            key={cat.value}
            type="button"
            onClick={() => setFilter(cat.value)}
            aria-pressed={filter === cat.value}
            style={{
              padding: "8px 16px",
              borderRadius: "20px",
              border: "1px solid #ccc",
              background: filter === cat.value ? "#1a237e" : "#fff",
              color: filter === cat.value ? "#fff" : "#333",
              cursor: "pointer",
              fontSize: "14px",
            }}
          >
            {cat.label}
          </button>
        ))}
      </nav>

      {loading && <p>読み込み中...</p>}

      {error && (
        <p role="alert" style={{ color: "#d32f2f" }}>
          {error}
        </p>
      )}

      {!loading && !error && (
        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(auto-fill, minmax(140px, 1fr))",
            gap: "16px",
          }}
        >
          {cards.map((card) => {
            const isCollected = collectedCardIds.has(card.id);
            return (
              <div
                key={card.id}
                style={{ opacity: isCollected ? 1 : 0.4 }}
              >
                <CollectionCard card={card} onSelect={handleSelectCard} />
              </div>
            );
          })}
        </div>
      )}

      {!loading && !error && cards.length === 0 && (
        <p>カードが見つかりませんでした。</p>
      )}

      {selectedCard && (
        <CardDetail card={selectedCard} onClose={handleCloseDetail} />
      )}
    </main>
  );
}
