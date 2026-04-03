import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { MatchmakingSpinner } from "./MatchmakingSpinner";

describe("MatchmakingSpinner", () => {
  it("renders default message", () => {
    render(<MatchmakingSpinner />);
    expect(screen.getByTestId("matchmaking-message")).toHaveTextContent(
      "対戦相手を探しています...",
    );
  });

  it("renders custom message", () => {
    render(<MatchmakingSpinner message="対戦相手を待っています..." />);
    expect(screen.getByTestId("matchmaking-message")).toHaveTextContent(
      "対戦相手を待っています...",
    );
  });

  it("renders spinner element", () => {
    render(<MatchmakingSpinner />);
    expect(screen.getByTestId("spinner")).toBeInTheDocument();
  });

  it("has accessible status role", () => {
    render(<MatchmakingSpinner />);
    expect(screen.getByRole("status")).toBeInTheDocument();
  });
});
