interface PictureCardProps {
  id: string;
  name: string;
  imageUrl: string;
  onGrab: (cardId: string) => void;
  disabled?: boolean;
  highlighted?: boolean;
}

export function PictureCard({
  id,
  name,
  imageUrl,
  onGrab,
  disabled = false,
  highlighted = false,
}: PictureCardProps) {
  return (
    <button
      type="button"
      onClick={() => onGrab(id)}
      disabled={disabled}
      aria-label={`取り札: ${name}`}
      data-testid={`picture-card-${id}`}
      style={{
        border: highlighted ? "3px solid #ffd700" : "2px solid #555",
        borderRadius: "12px",
        padding: "8px",
        background: highlighted
          ? "rgba(255, 215, 0, 0.15)"
          : "rgba(255, 255, 255, 0.05)",
        cursor: disabled ? "not-allowed" : "pointer",
        opacity: disabled ? 0.5 : 1,
        textAlign: "center",
        transition: "transform 0.15s ease, border-color 0.15s ease",
        width: "100%",
      }}
    >
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={imageUrl}
        alt={name}
        style={{
          width: "100%",
          height: "120px",
          objectFit: "cover",
          borderRadius: "8px",
        }}
      />
      <p
        style={{
          marginTop: "8px",
          fontSize: "14px",
          fontWeight: "bold",
          color: "#fff",
        }}
      >
        {name}
      </p>
    </button>
  );
}
