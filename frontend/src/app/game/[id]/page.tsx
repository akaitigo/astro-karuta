"use client";

import { useCallback, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { ReadingCard } from "@/components/ReadingCard";
import { PictureCard } from "@/components/PictureCard";
import { ScoreBoard } from "@/components/ScoreBoard";
import { GameResult } from "@/components/GameResult";
import { MatchmakingSpinner } from "@/components/MatchmakingSpinner";
import { useGame } from "@/hooks/useGame";
import { useAudio } from "@/hooks/useAudio";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/ws";

export default function GamePage() {
  const params = useParams();
  const router = useRouter();
  const gameId = typeof params.id === "string" ? params.id : "";

  const { state, connectionStatus, grabCard, reset, players } = useGame({
    wsUrl: WS_URL,
  });

  const { speak } = useAudio();

  // Auto-read reading text aloud when a new card is revealed
  useEffect(() => {
    if (state.readingText && state.gameStatus === "playing") {
      speak(state.readingText);
    }
  }, [state.readingText, state.gameStatus, speak]);

  const handleGrab = useCallback(
    (cardId: string) => {
      if (state.lastGrabResult !== null) return; // Already grabbed this round
      grabCard(cardId);
    },
    [grabCard, state.lastGrabResult],
  );

  const handleBackToLobby = useCallback(() => {
    reset();
    router.push("/lobby");
  }, [reset, router]);

  // Connection status
  if (connectionStatus === "disconnected" && state.gameStatus === "lobby") {
    return (
      <main
        style={{
          minHeight: "100vh",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          gap: "16px",
          background: "#0a0a1a",
          color: "#fff",
        }}
      >
        <p>ゲーム: {gameId}</p>
        <p style={{ color: "#aaa" }}>
          接続が切断されました。ロビーに戻ってください。
        </p>
        <button
          type="button"
          onClick={handleBackToLobby}
          style={{
            padding: "12px 24px",
            fontSize: "16px",
            color: "#fff",
            background: "#6495ed",
            border: "none",
            borderRadius: "8px",
            cursor: "pointer",
          }}
        >
          ロビーに戻る
        </button>
      </main>
    );
  }

  // Waiting for opponent
  if (
    state.gameStatus === "waiting" ||
    state.gameStatus === "matchmaking"
  ) {
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
        <MatchmakingSpinner message="対戦相手を待っています..." />
      </main>
    );
  }

  // Game over
  if (state.gameStatus === "finished" && state.gameResult) {
    return (
      <main
        style={{
          minHeight: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          background: "#0a0a1a",
          color: "#fff",
        }}
      >
        <GameResult
          result={state.gameResult}
          currentPlayerId={state.playerId}
          onBackToLobby={handleBackToLobby}
        />
      </main>
    );
  }

  // Active game
  const isGrabbed = state.lastGrabResult !== null;

  return (
    <main
      style={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        background: "#0a0a1a",
        color: "#fff",
        padding: "16px",
      }}
    >
      {/* Score Board */}
      <ScoreBoard players={players} currentPlayerId={state.playerId} />

      {/* Remaining cards */}
      <p
        style={{
          textAlign: "center",
          fontSize: "12px",
          color: "#888",
          margin: "8px 0",
        }}
        data-testid="remaining-cards"
      >
        残り {state.totalCards - state.cardIndex} / {state.totalCards} 枚
      </p>

      {/* Reading Card */}
      <div style={{ margin: "16px 0" }}>
        <ReadingCard
          readingText={state.readingText}
          cardIndex={state.cardIndex}
          totalCards={state.totalCards}
        />
      </div>

      {/* Grab result feedback */}
      {isGrabbed && state.lastGrabResult && (
        <div
          style={{
            textAlign: "center",
            padding: "8px",
            marginBottom: "8px",
          }}
          data-testid="grab-feedback"
        >
          {state.lastGrabResult.correct ? (
            <p style={{ color: "#4caf50", fontWeight: "bold" }}>
              {state.lastGrabResult.winner_name} が正解!
              ({state.lastGrabResult.card_name})
            </p>
          ) : (
            <p style={{ color: "#ff6b6b", fontWeight: "bold" }}>
              お手つき!
            </p>
          )}
        </div>
      )}

      {/* Picture Card Grid */}
      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(2, 1fr)",
          gap: "12px",
          maxWidth: "480px",
          margin: "0 auto",
          width: "100%",
        }}
        data-testid="card-grid"
      >
        {state.candidates.map((card) => (
          <PictureCard
            key={card.id}
            id={card.id}
            name={card.name}
            imageUrl={card.image_url}
            onGrab={handleGrab}
            disabled={isGrabbed}
            highlighted={
              isGrabbed &&
              state.lastGrabResult !== null &&
              state.lastGrabResult.card_id === card.id
            }
          />
        ))}
      </div>

      {/* Error */}
      {state.error && (
        <p
          role="alert"
          style={{
            color: "#ff6b6b",
            textAlign: "center",
            marginTop: "16px",
          }}
        >
          {state.error}
        </p>
      )}
    </main>
  );
}
