import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { ReadingCard } from "./ReadingCard";

describe("ReadingCard", () => {
  it("renders the reading text", () => {
    render(
      <ReadingCard
        readingText="おおぐま座の北斗七星"
        cardIndex={3}
        totalCards={20}
      />,
    );
    expect(screen.getByTestId("reading-text")).toHaveTextContent(
      "おおぐま座の北斗七星",
    );
  });

  it("renders the card progress", () => {
    render(
      <ReadingCard readingText="test" cardIndex={5} totalCards={20} />,
    );
    expect(screen.getByText("5 / 20")).toBeInTheDocument();
  });

  it("has accessible region label", () => {
    render(
      <ReadingCard readingText="test" cardIndex={1} totalCards={10} />,
    );
    expect(screen.getByRole("region", { name: "読み札" })).toBeInTheDocument();
  });
});
