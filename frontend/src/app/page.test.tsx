import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import Home from "./page";

describe("Home", () => {
  it("renders the title", () => {
    render(<Home />);
    expect(screen.getByText("Astro-Karuta")).toBeInTheDocument();
  });

  it("renders the description", () => {
    render(<Home />);
    expect(
      screen.getByText("天文知識をかるた形式で学ぼう"),
    ).toBeInTheDocument();
  });
});
