import type { PlayerResult } from "@/types/game";

interface ScoreBoardProps {
  players: PlayerResult[];
  currentPlayerId: string;
}

export function ScoreBoard({ players, currentPlayerId }: ScoreBoardProps) {
  return (
    <div
      role="region"
      aria-label="スコアボード"
      style={{
        display: "flex",
        justifyContent: "center",
        gap: "24px",
        padding: "12px",
      }}
    >
      {players.map((player) => {
        const isMe = player.player_id === currentPlayerId;
        return (
          <div
            key={player.player_id}
            data-testid={`score-${player.player_id}`}
            style={{
              padding: "8px 16px",
              borderRadius: "8px",
              background: isMe
                ? "rgba(100, 149, 237, 0.3)"
                : "rgba(255, 255, 255, 0.1)",
              border: isMe ? "1px solid #6495ed" : "1px solid #444",
              textAlign: "center",
              minWidth: "100px",
            }}
          >
            <p
              style={{
                fontSize: "12px",
                color: "#aaa",
                marginBottom: "4px",
              }}
            >
              {isMe ? "あなた" : player.player_name}
            </p>
            <p
              style={{
                fontSize: "28px",
                fontWeight: "bold",
                color: "#fff",
              }}
            >
              {player.score}
            </p>
          </div>
        );
      })}
    </div>
  );
}
