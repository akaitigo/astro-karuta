interface MatchmakingSpinnerProps {
  message?: string;
}

export function MatchmakingSpinner({
  message = "対戦相手を探しています...",
}: MatchmakingSpinnerProps) {
  return (
    <div
      role="status"
      aria-label="マッチメイキング中"
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        gap: "24px",
        padding: "48px",
      }}
    >
      <div
        data-testid="spinner"
        style={{
          width: "48px",
          height: "48px",
          border: "4px solid rgba(255, 255, 255, 0.2)",
          borderTop: "4px solid #ffd700",
          borderRadius: "50%",
          animation: "spin 1s linear infinite",
        }}
      />
      <p
        style={{
          fontSize: "16px",
          color: "#ccc",
        }}
        data-testid="matchmaking-message"
      >
        {message}
      </p>
      <style>{`
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}
