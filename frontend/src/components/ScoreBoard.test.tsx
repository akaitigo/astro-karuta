import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { ScoreBoard } from "./ScoreBoard";
import type { PlayerResult } from "@/types/game";

describe("ScoreBoard", () => {
  const players: PlayerResult[] = [
    {
      player_id: "p1",
      player_name: "Alice",
      score: 5,
      captured_ids: [],
    },
    {
      player_id: "p2",
      player_name: "Bob",
      score: 3,
      captured_ids: [],
    },
  ];

  it("renders both players", () => {
    render(<ScoreBoard players={players} currentPlayerId="p1" />);
    expect(screen.getByTestId("score-p1")).toBeInTheDocument();
    expect(screen.getByTestId("score-p2")).toBeInTheDocument();
  });

  it("displays scores", () => {
    render(<ScoreBoard players={players} currentPlayerId="p1" />);
    expect(screen.getByTestId("score-p1")).toHaveTextContent("5");
    expect(screen.getByTestId("score-p2")).toHaveTextContent("3");
  });

  it("shows 'あなた' for current player", () => {
    render(<ScoreBoard players={players} currentPlayerId="p1" />);
    expect(screen.getByText("あなた")).toBeInTheDocument();
  });

  it("shows opponent name", () => {
    render(<ScoreBoard players={players} currentPlayerId="p1" />);
    expect(screen.getByText("Bob")).toBeInTheDocument();
  });

  it("has accessible region label", () => {
    render(<ScoreBoard players={players} currentPlayerId="p1" />);
    expect(
      screen.getByRole("region", { name: "スコアボード" }),
    ).toBeInTheDocument();
  });
});
