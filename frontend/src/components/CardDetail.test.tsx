import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";
import CardDetail from "./CardDetail";
import type { Card } from "@/types/card";

const mockCard: Card = {
  id: "card-1",
  name: "オリオン座",
  category: "constellation",
  reading_text: "冬の夜空に三つ星が並ぶ、狩人の姿",
  image_url: "/images/orion.jpg",
  description: "冬を代表する星座。ベテルギウスとリゲルが目印。",
  magnitude: 0.12,
  distance: "1,344光年",
  best_season: "winter",
};

describe("CardDetail", () => {
  it("renders the card name as heading", () => {
    render(<CardDetail card={mockCard} onClose={vi.fn()} />);
    expect(
      screen.getByRole("heading", { name: "オリオン座" })
    ).toBeInTheDocument();
  });

  it("renders the reading text", () => {
    render(<CardDetail card={mockCard} onClose={vi.fn()} />);
    expect(
      screen.getByText("冬の夜空に三つ星が並ぶ、狩人の姿")
    ).toBeInTheDocument();
  });

  it("renders the description", () => {
    render(<CardDetail card={mockCard} onClose={vi.fn()} />);
    expect(
      screen.getByText("冬を代表する星座。ベテルギウスとリゲルが目印。")
    ).toBeInTheDocument();
  });

  it("renders the category label", () => {
    render(<CardDetail card={mockCard} onClose={vi.fn()} />);
    expect(screen.getByText("星座")).toBeInTheDocument();
  });

  it("renders the season label", () => {
    render(<CardDetail card={mockCard} onClose={vi.fn()} />);
    expect(screen.getByText("冬")).toBeInTheDocument();
  });

  it("renders magnitude when present", () => {
    render(<CardDetail card={mockCard} onClose={vi.fn()} />);
    expect(screen.getByText("等級")).toBeInTheDocument();
    expect(screen.getByText("0.12")).toBeInTheDocument();
  });

  it("renders distance when present", () => {
    render(<CardDetail card={mockCard} onClose={vi.fn()} />);
    expect(screen.getByText("距離")).toBeInTheDocument();
    expect(screen.getByText("1,344光年")).toBeInTheDocument();
  });

  it("does not render magnitude when absent", () => {
    const cardWithoutMagnitude: Card = {
      ...mockCard,
      magnitude: undefined,
    };
    render(<CardDetail card={cardWithoutMagnitude} onClose={vi.fn()} />);
    expect(screen.queryByText("等級")).not.toBeInTheDocument();
  });

  it("does not render distance when absent", () => {
    const cardWithoutDistance: Card = {
      ...mockCard,
      distance: undefined,
    };
    render(<CardDetail card={cardWithoutDistance} onClose={vi.fn()} />);
    expect(screen.queryByText("距離")).not.toBeInTheDocument();
  });

  it("calls onClose when close button is clicked", async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();
    render(<CardDetail card={mockCard} onClose={onClose} />);

    await user.click(screen.getByLabelText("閉じる"));
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it("has dialog role with accessible label", () => {
    render(<CardDetail card={mockCard} onClose={vi.fn()} />);
    expect(
      screen.getByRole("dialog", { name: "オリオン座の詳細" })
    ).toBeInTheDocument();
  });

  it("renders all-season label correctly", () => {
    const allSeasonCard: Card = {
      ...mockCard,
      best_season: "all",
    };
    render(<CardDetail card={allSeasonCard} onClose={vi.fn()} />);
    expect(screen.getByText("通年")).toBeInTheDocument();
  });
});
