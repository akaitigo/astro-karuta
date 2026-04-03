"use client";

import type { Card } from "@/types/card";

interface CardDetailProps {
  card: Card;
  onClose: () => void;
}

const categoryLabels: Record<string, string> = {
  constellation: "星座",
  planet: "惑星",
  phenomenon: "天文現象",
};

const seasonLabels: Record<string, string> = {
  spring: "春",
  summer: "夏",
  autumn: "秋",
  winter: "冬",
  all: "通年",
};

export default function CardDetail({ card, onClose }: CardDetailProps) {
  const categoryLabel = categoryLabels[card.category] ?? card.category;
  const seasonLabel = seasonLabels[card.best_season] ?? card.best_season;

  return (
    <div
      role="dialog"
      aria-label={`${card.name}の詳細`}
      style={{
        position: "fixed",
        inset: 0,
        background: "rgba(0, 0, 0, 0.6)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        zIndex: 100,
        padding: "16px",
      }}
    >
      <div
        style={{
          background: "#fff",
          borderRadius: "12px",
          padding: "24px",
          maxWidth: "480px",
          width: "100%",
          maxHeight: "90vh",
          overflow: "auto",
        }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            marginBottom: "16px",
          }}
        >
          <h2 style={{ margin: 0, fontSize: "20px" }}>{card.name}</h2>
          <button
            type="button"
            onClick={onClose}
            aria-label="閉じる"
            style={{
              background: "none",
              border: "none",
              fontSize: "24px",
              cursor: "pointer",
              padding: "4px 8px",
            }}
          >
            ×
          </button>
        </div>

        <p
          style={{
            fontStyle: "italic",
            color: "#555",
            marginBottom: "12px",
            lineHeight: 1.6,
          }}
        >
          {card.reading_text}
        </p>

        <p style={{ lineHeight: 1.6, marginBottom: "16px" }}>
          {card.description}
        </p>

        <dl
          style={{
            display: "grid",
            gridTemplateColumns: "auto 1fr",
            gap: "8px 16px",
            margin: 0,
          }}
        >
          <dt style={{ fontWeight: "bold" }}>カテゴリ</dt>
          <dd style={{ margin: 0 }}>{categoryLabel}</dd>

          <dt style={{ fontWeight: "bold" }}>見頃</dt>
          <dd style={{ margin: 0 }}>{seasonLabel}</dd>

          {card.magnitude !== undefined && (
            <>
              <dt style={{ fontWeight: "bold" }}>等級</dt>
              <dd style={{ margin: 0 }}>{card.magnitude}</dd>
            </>
          )}

          {card.distance !== undefined && (
            <>
              <dt style={{ fontWeight: "bold" }}>距離</dt>
              <dd style={{ margin: 0 }}>{card.distance}</dd>
            </>
          )}
        </dl>
      </div>
    </div>
  );
}
