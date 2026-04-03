import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect, vi } from "vitest";
import CollectionCard from "./CollectionCard";
import type { Card } from "@/types/card";

const mockCard: Card = {
  id: "card-1",
  name: "オリオン座",
  category: "constellation",
  reading_text: "冬の夜空に三つ星が並ぶ、狩人の姿",
  image_url: "/images/orion.jpg",
  description: "冬を代表する星座。",
  magnitude: 0.12,
  distance: "1,344光年",
  best_season: "winter",
};

describe("CollectionCard", () => {
  it("renders the card name", () => {
    render(<CollectionCard card={mockCard} onSelect={vi.fn()} />);
    expect(screen.getByText("オリオン座")).toBeInTheDocument();
  });

  it("renders the category label", () => {
    render(<CollectionCard card={mockCard} onSelect={vi.fn()} />);
    expect(screen.getByText("星座")).toBeInTheDocument();
  });

  it("renders planet category label", () => {
    const planetCard: Card = {
      ...mockCard,
      id: "card-2",
      name: "火星",
      category: "planet",
    };
    render(<CollectionCard card={planetCard} onSelect={vi.fn()} />);
    expect(screen.getByText("惑星")).toBeInTheDocument();
  });

  it("renders phenomenon category label", () => {
    const phenomenonCard: Card = {
      ...mockCard,
      id: "card-3",
      name: "皆既日食",
      category: "phenomenon",
    };
    render(<CollectionCard card={phenomenonCard} onSelect={vi.fn()} />);
    expect(screen.getByText("天文現象")).toBeInTheDocument();
  });

  it("calls onSelect when clicked", async () => {
    const user = userEvent.setup();
    const onSelect = vi.fn();
    render(<CollectionCard card={mockCard} onSelect={onSelect} />);

    await user.click(screen.getByRole("button"));
    expect(onSelect).toHaveBeenCalledTimes(1);
    expect(onSelect).toHaveBeenCalledWith(mockCard);
  });

  it("has accessible label", () => {
    render(<CollectionCard card={mockCard} onSelect={vi.fn()} />);
    expect(
      screen.getByLabelText("オリオン座の詳細を表示")
    ).toBeInTheDocument();
  });

  it("displays the first character as avatar", () => {
    render(<CollectionCard card={mockCard} onSelect={vi.fn()} />);
    expect(screen.getByText("オ")).toBeInTheDocument();
  });
});
