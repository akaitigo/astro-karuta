"use client";

import type { Card } from "@/types/card";

interface CollectionCardProps {
  card: Card;
  onSelect: (card: Card) => void;
}

const categoryLabels: Record<string, string> = {
  constellation: "星座",
  planet: "惑星",
  phenomenon: "天文現象",
};

const categoryColors: Record<string, string> = {
  constellation: "#1a237e",
  planet: "#b71c1c",
  phenomenon: "#1b5e20",
};

export default function CollectionCard({ card, onSelect }: CollectionCardProps) {
  const label = categoryLabels[card.category] ?? card.category;
  const borderColor = categoryColors[card.category] ?? "#333";

  return (
    <button
      type="button"
      onClick={() => onSelect(card)}
      aria-label={`${card.name}の詳細を表示`}
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        padding: "12px",
        border: `2px solid ${borderColor}`,
        borderRadius: "8px",
        background: "#fff",
        cursor: "pointer",
        width: "100%",
        textAlign: "center",
      }}
    >
      <div
        style={{
          width: "80px",
          height: "80px",
          borderRadius: "50%",
          background: `linear-gradient(135deg, ${borderColor}, #555)`,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          color: "#fff",
          fontSize: "28px",
          marginBottom: "8px",
        }}
        aria-hidden="true"
      >
        {card.name.charAt(0)}
      </div>
      <span style={{ fontWeight: "bold", fontSize: "14px" }}>{card.name}</span>
      <span
        style={{
          fontSize: "11px",
          color: borderColor,
          marginTop: "4px",
          padding: "2px 8px",
          border: `1px solid ${borderColor}`,
          borderRadius: "12px",
        }}
      >
        {label}
      </span>
    </button>
  );
}
