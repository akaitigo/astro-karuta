import type { GameOverPayload } from "@/types/game";

interface GameResultProps {
  result: GameOverPayload;
  currentPlayerId: string;
  onBackToLobby: () => void;
}

export function GameResult({
  result,
  currentPlayerId,
  onBackToLobby,
}: GameResultProps) {
  const isWinner = result.winner_id === currentPlayerId;
  const sortedPlayers = [...result.players].sort(
    (a, b) => b.score - a.score,
  );

  return (
    <div
      role="region"
      aria-label="対戦結果"
      style={{
        textAlign: "center",
        padding: "32px",
        maxWidth: "480px",
        margin: "0 auto",
      }}
    >
      <h2
        style={{
          fontSize: "32px",
          fontWeight: "bold",
          color: isWinner ? "#ffd700" : "#6495ed",
          marginBottom: "16px",
        }}
        data-testid="result-heading"
      >
        {isWinner ? "勝利!" : "惜敗..."}
      </h2>

      <div style={{ marginBottom: "24px" }}>
        {sortedPlayers.map((player, index) => (
          <div
            key={player.player_id}
            data-testid={`result-player-${player.player_id}`}
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              padding: "12px 16px",
              marginBottom: "8px",
              borderRadius: "8px",
              background:
                index === 0
                  ? "rgba(255, 215, 0, 0.15)"
                  : "rgba(255, 255, 255, 0.05)",
              border:
                index === 0 ? "1px solid #ffd700" : "1px solid #444",
            }}
          >
            <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
              <span
                style={{
                  fontSize: "20px",
                  fontWeight: "bold",
                  color: index === 0 ? "#ffd700" : "#888",
                  width: "32px",
                }}
              >
                {index + 1}.
              </span>
              <span style={{ color: "#fff", fontSize: "16px" }}>
                {player.player_id === currentPlayerId
                  ? "あなた"
                  : player.player_name}
              </span>
            </div>
            <span
              style={{
                fontSize: "24px",
                fontWeight: "bold",
                color: "#fff",
              }}
            >
              {player.score} 枚
            </span>
          </div>
        ))}
      </div>

      <button
        type="button"
        onClick={onBackToLobby}
        data-testid="back-to-lobby"
        style={{
          padding: "12px 32px",
          fontSize: "16px",
          fontWeight: "bold",
          color: "#fff",
          background: "#6495ed",
          border: "none",
          borderRadius: "8px",
          cursor: "pointer",
        }}
      >
        ロビーに戻る
      </button>
    </div>
  );
}
