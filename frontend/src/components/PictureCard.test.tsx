import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { PictureCard } from "./PictureCard";

describe("PictureCard", () => {
  const defaultProps = {
    id: "card-1",
    name: "オリオン座",
    imageUrl: "/images/orion.png",
    onGrab: vi.fn(),
  };

  it("renders the card name", () => {
    render(<PictureCard {...defaultProps} />);
    expect(screen.getByText("オリオン座")).toBeInTheDocument();
  });

  it("renders the card image with correct alt text", () => {
    render(<PictureCard {...defaultProps} />);
    const img = screen.getByAltText("オリオン座");
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute("src", "/images/orion.png");
  });

  it("calls onGrab when clicked", () => {
    const onGrab = vi.fn();
    render(<PictureCard {...defaultProps} onGrab={onGrab} />);
    fireEvent.click(screen.getByRole("button"));
    expect(onGrab).toHaveBeenCalledWith("card-1");
  });

  it("does not call onGrab when disabled", () => {
    const onGrab = vi.fn();
    render(<PictureCard {...defaultProps} onGrab={onGrab} disabled />);
    fireEvent.click(screen.getByRole("button"));
    expect(onGrab).not.toHaveBeenCalled();
  });

  it("has accessible aria-label", () => {
    render(<PictureCard {...defaultProps} />);
    expect(
      screen.getByRole("button", { name: "取り札: オリオン座" }),
    ).toBeInTheDocument();
  });

  it("applies highlighted style", () => {
    render(<PictureCard {...defaultProps} highlighted />);
    const button = screen.getByRole("button");
    // jsdom normalizes hex colors to rgb
    expect(button.style.border).toContain("solid");
    expect(button.style.border).toContain("3px");
  });
});
