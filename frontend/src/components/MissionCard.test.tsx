import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import MissionCard from "./MissionCard";
import type { UserMission } from "@/types/mission";

function createMission(overrides: Partial<UserMission> = {}): UserMission {
  const now = new Date();
  const validTo = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000);

  return {
    id: "mission-1",
    user_id: "user-1",
    mission_id: "obs-7-card-1",
    card_id: "card-1",
    title: "オリオン座を観測しよう",
    description: "夜空でオリオン座を見つけて、観測ボタンを押そう！",
    status: "active",
    valid_from: now.toISOString(),
    valid_to: validTo.toISOString(),
    created_at: now.toISOString(),
    ...overrides,
  };
}

describe("MissionCard", () => {
  it("renders mission title", () => {
    const mission = createMission();
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    expect(screen.getByText("オリオン座を観測しよう")).toBeInTheDocument();
  });

  it("renders mission description", () => {
    const mission = createMission();
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    expect(screen.getByText("夜空でオリオン座を見つけて、観測ボタンを押そう！")).toBeInTheDocument();
  });

  it("shows active status for active missions", () => {
    const mission = createMission({ status: "active" });
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    expect(screen.getByTestId("mission-status")).toHaveTextContent("進行中");
  });

  it("shows completed status for completed missions", () => {
    const mission = createMission({
      status: "completed",
      completed_at: new Date().toISOString(),
    });
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    expect(screen.getByTestId("mission-status")).toHaveTextContent("達成！");
  });

  it("shows expired status for expired missions", () => {
    const mission = createMission({ status: "expired" });
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    expect(screen.getByTestId("mission-status")).toHaveTextContent("期限切れ");
  });

  it("shows complete button for active missions", () => {
    const mission = createMission({ status: "active" });
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    expect(screen.getByTestId("mission-complete-button")).toBeInTheDocument();
    expect(screen.getByTestId("mission-complete-button")).toHaveTextContent("観測した！");
  });

  it("hides complete button for completed missions", () => {
    const mission = createMission({ status: "completed" });
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    expect(screen.queryByTestId("mission-complete-button")).not.toBeInTheDocument();
  });

  it("hides complete button for expired missions", () => {
    const mission = createMission({ status: "expired" });
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    expect(screen.queryByTestId("mission-complete-button")).not.toBeInTheDocument();
  });

  it("calls onComplete when button is clicked", () => {
    const mission = createMission();
    const onComplete = vi.fn();
    render(<MissionCard mission={mission} onComplete={onComplete} />);

    fireEvent.click(screen.getByTestId("mission-complete-button"));
    expect(onComplete).toHaveBeenCalledWith("mission-1");
  });

  it("shows loading state", () => {
    const mission = createMission();
    render(<MissionCard mission={mission} onComplete={() => {}} isLoading={true} />);
    const button = screen.getByTestId("mission-complete-button");
    expect(button).toHaveTextContent("送信中...");
    expect(button).toBeDisabled();
  });

  it("shows remaining time for active missions", () => {
    const mission = createMission();
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    const remaining = screen.getByTestId("mission-remaining");
    expect(remaining.textContent).toMatch(/残り/);
  });

  it("shows observation complete for completed missions", () => {
    const mission = createMission({ status: "completed" });
    render(<MissionCard mission={mission} onComplete={() => {}} />);
    expect(screen.getByTestId("mission-remaining")).toHaveTextContent("観測完了");
  });
});
