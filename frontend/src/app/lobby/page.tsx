"use client";

import { useCallback, useState } from "react";
import { useRouter } from "next/navigation";
import { MatchmakingSpinner } from "@/components/MatchmakingSpinner";
import { useWebSocket } from "@/hooks/useWebSocket";
import type { WSMessage } from "@/types/game";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/ws";

export default function LobbyPage() {
  const router = useRouter();
  const [playerName, setPlayerName] = useState("");
  const [roomCode, setRoomCode] = useState("");
  const [isMatchmaking, setIsMatchmaking] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleMessage = useCallback(
    (message: WSMessage) => {
      switch (message.type) {
        case "player_joined": {
          if (
            typeof message.payload === "object" &&
            message.payload !== null &&
            "room_code" in message.payload
          ) {
            const payload = message.payload as { room_code: string };
            router.push(`/game/${encodeURIComponent(payload.room_code)}`);
          }
          break;
        }
        case "match_found": {
          if (
            typeof message.payload === "object" &&
            message.payload !== null &&
            "room_code" in message.payload
          ) {
            const payload = message.payload as { room_code: string };
            router.push(`/game/${encodeURIComponent(payload.room_code)}`);
          }
          break;
        }
        case "waiting":
          setIsMatchmaking(true);
          break;
        case "error": {
          if (
            typeof message.payload === "object" &&
            message.payload !== null &&
            "message" in message.payload
          ) {
            const payload = message.payload as { message: string };
            setError(payload.message);
          }
          setIsMatchmaking(false);
          break;
        }
      }
    },
    [router],
  );

  const { send, connect, status } = useWebSocket({
    url: WS_URL,
    onMessage: handleMessage,
    autoConnect: false,
  });

  const ensureConnectedAndSend = useCallback(
    (
      roomCodeVal: string,
      playerNameVal: string,
      randomMatch: boolean,
    ) => {
      setError(null);
      if (!playerNameVal.trim()) {
        setError("プレイヤー名を入力してください");
        return;
      }

      connect();
      // Wait for connection to establish before sending
      setTimeout(() => {
        send("join", {
          room_code: roomCodeVal,
          player_name: playerNameVal,
          random_match: randomMatch,
        });
      }, 200);
    },
    [connect, send],
  );

  const handleRandomMatch = useCallback(() => {
    setIsMatchmaking(true);
    ensureConnectedAndSend("", playerName, true);
  }, [playerName, ensureConnectedAndSend]);

  const handleCreateRoom = useCallback(() => {
    ensureConnectedAndSend("", playerName, false);
  }, [playerName, ensureConnectedAndSend]);

  const handleJoinRoom = useCallback(() => {
    if (!roomCode.trim()) {
      setError("ルームコードを入力してください");
      return;
    }
    ensureConnectedAndSend(roomCode, playerName, false);
  }, [roomCode, playerName, ensureConnectedAndSend]);

  if (isMatchmaking) {
    return (
      <main
        style={{
          minHeight: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          background: "#0a0a1a",
        }}
      >
        <MatchmakingSpinner />
      </main>
    );
  }

  return (
    <main
      style={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        gap: "24px",
        padding: "24px",
        background: "#0a0a1a",
        color: "#fff",
      }}
    >
      <h1 style={{ fontSize: "32px", fontWeight: "bold", color: "#ffd700" }}>
        Astro-Karuta
      </h1>
      <p style={{ color: "#aaa", fontSize: "14px" }}>天文かるた対戦ロビー</p>

      {/* Player name input */}
      <div style={{ width: "100%", maxWidth: "320px" }}>
        <label
          htmlFor="player-name"
          style={{ display: "block", marginBottom: "4px", fontSize: "14px" }}
        >
          プレイヤー名
        </label>
        <input
          id="player-name"
          type="text"
          value={playerName}
          onChange={(e) => setPlayerName(e.target.value)}
          placeholder="名前を入力"
          maxLength={20}
          style={{
            width: "100%",
            padding: "10px 12px",
            borderRadius: "8px",
            border: "1px solid #444",
            background: "rgba(255,255,255,0.1)",
            color: "#fff",
            fontSize: "16px",
          }}
        />
      </div>

      {/* Random Match */}
      <button
        type="button"
        onClick={handleRandomMatch}
        disabled={status === "connecting"}
        data-testid="random-match-btn"
        style={{
          width: "100%",
          maxWidth: "320px",
          padding: "14px",
          fontSize: "16px",
          fontWeight: "bold",
          color: "#0a0a1a",
          background: "#ffd700",
          border: "none",
          borderRadius: "8px",
          cursor: "pointer",
        }}
      >
        ランダムマッチ
      </button>

      {/* Divider */}
      <div
        style={{
          display: "flex",
          alignItems: "center",
          gap: "12px",
          width: "100%",
          maxWidth: "320px",
        }}
      >
        <hr style={{ flex: 1, border: "none", borderTop: "1px solid #444" }} />
        <span style={{ color: "#888", fontSize: "12px" }}>または</span>
        <hr style={{ flex: 1, border: "none", borderTop: "1px solid #444" }} />
      </div>

      {/* Room code input + Join */}
      <div style={{ width: "100%", maxWidth: "320px" }}>
        <label
          htmlFor="room-code"
          style={{ display: "block", marginBottom: "4px", fontSize: "14px" }}
        >
          ルームコード
        </label>
        <div style={{ display: "flex", gap: "8px" }}>
          <input
            id="room-code"
            type="text"
            value={roomCode}
            onChange={(e) => setRoomCode(e.target.value.toUpperCase())}
            placeholder="ABC123"
            maxLength={6}
            style={{
              flex: 1,
              padding: "10px 12px",
              borderRadius: "8px",
              border: "1px solid #444",
              background: "rgba(255,255,255,0.1)",
              color: "#fff",
              fontSize: "16px",
              letterSpacing: "2px",
              textTransform: "uppercase",
            }}
          />
          <button
            type="button"
            onClick={handleJoinRoom}
            disabled={status === "connecting"}
            data-testid="join-room-btn"
            style={{
              padding: "10px 20px",
              fontSize: "14px",
              fontWeight: "bold",
              color: "#fff",
              background: "#6495ed",
              border: "none",
              borderRadius: "8px",
              cursor: "pointer",
              whiteSpace: "nowrap",
            }}
          >
            参加
          </button>
        </div>
      </div>

      {/* Create Room */}
      <button
        type="button"
        onClick={handleCreateRoom}
        disabled={status === "connecting"}
        data-testid="create-room-btn"
        style={{
          width: "100%",
          maxWidth: "320px",
          padding: "14px",
          fontSize: "16px",
          fontWeight: "bold",
          color: "#fff",
          background: "transparent",
          border: "2px solid #6495ed",
          borderRadius: "8px",
          cursor: "pointer",
        }}
      >
        ルームを作成
      </button>

      {/* Error display */}
      {error && (
        <p
          role="alert"
          style={{ color: "#ff6b6b", fontSize: "14px" }}
          data-testid="lobby-error"
        >
          {error}
        </p>
      )}
    </main>
  );
}
