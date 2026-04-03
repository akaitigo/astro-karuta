interface ReadingCardProps {
  readingText: string;
  cardIndex: number;
  totalCards: number;
}

export function ReadingCard({
  readingText,
  cardIndex,
  totalCards,
}: ReadingCardProps) {
  return (
    <div
      role="region"
      aria-label="読み札"
      style={{
        background: "linear-gradient(135deg, #1a1a3e 0%, #2d1b69 100%)",
        border: "2px solid #ffd700",
        borderRadius: "16px",
        padding: "32px",
        textAlign: "center",
        maxWidth: "480px",
        margin: "0 auto",
      }}
    >
      <p
        style={{
          fontSize: "12px",
          color: "#aaa",
          marginBottom: "8px",
        }}
      >
        {cardIndex} / {totalCards}
      </p>
      <p
        style={{
          fontSize: "24px",
          fontWeight: "bold",
          color: "#fff",
          lineHeight: 1.6,
        }}
        data-testid="reading-text"
      >
        {readingText}
      </p>
    </div>
  );
}
