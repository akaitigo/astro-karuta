import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { GameResult } from "./GameResult";
import type { GameOverPayload } from "@/types/game";

describe("GameResult", () => {
  const result: GameOverPayload = {
    players: [
      {
        player_id: "p1",
        player_name: "Alice",
        score: 8,
        captured_ids: ["c1", "c2"],
      },
      {
        player_id: "p2",
        player_name: "Bob",
        score: 5,
        captured_ids: ["c3"],
      },
    ],
    winner_id: "p1",
  };

  it("shows victory message for winner", () => {
    render(
      <GameResult
        result={result}
        currentPlayerId="p1"
        onBackToLobby={vi.fn()}
      />,
    );
    expect(screen.getByTestId("result-heading")).toHaveTextContent("勝利!");
  });

  it("shows defeat message for loser", () => {
    render(
      <GameResult
        result={result}
        currentPlayerId="p2"
        onBackToLobby={vi.fn()}
      />,
    );
    expect(screen.getByTestId("result-heading")).toHaveTextContent("惜敗...");
  });

  it("renders player scores sorted by score", () => {
    render(
      <GameResult
        result={result}
        currentPlayerId="p1"
        onBackToLobby={vi.fn()}
      />,
    );
    const p1 = screen.getByTestId("result-player-p1");
    const p2 = screen.getByTestId("result-player-p2");
    expect(p1).toHaveTextContent("8 枚");
    expect(p2).toHaveTextContent("5 枚");
  });

  it("shows 'あなた' for current player", () => {
    render(
      <GameResult
        result={result}
        currentPlayerId="p1"
        onBackToLobby={vi.fn()}
      />,
    );
    expect(screen.getByText("あなた")).toBeInTheDocument();
  });

  it("calls onBackToLobby when button clicked", () => {
    const onBackToLobby = vi.fn();
    render(
      <GameResult
        result={result}
        currentPlayerId="p1"
        onBackToLobby={onBackToLobby}
      />,
    );
    fireEvent.click(screen.getByTestId("back-to-lobby"));
    expect(onBackToLobby).toHaveBeenCalledOnce();
  });

  it("has accessible region label", () => {
    render(
      <GameResult
        result={result}
        currentPlayerId="p1"
        onBackToLobby={vi.fn()}
      />,
    );
    expect(
      screen.getByRole("region", { name: "対戦結果" }),
    ).toBeInTheDocument();
  });
});
