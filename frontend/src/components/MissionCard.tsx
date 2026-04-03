"use client";

import type { UserMission } from "@/types/mission";

interface MissionCardProps {
  mission: UserMission;
  onComplete: (missionId: string) => void;
  isLoading?: boolean;
}

function formatRemainingTime(validTo: string): string {
  const end = new Date(validTo);
  const now = new Date();
  const diffMs = end.getTime() - now.getTime();

  if (diffMs <= 0) {
    return "期限切れ";
  }

  const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  const hours = Math.floor((diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));

  if (days > 0) {
    return `残り${days}日${hours}時間`;
  }
  return `残り${hours}時間`;
}

function statusLabel(status: UserMission["status"]): string {
  switch (status) {
    case "active":
      return "進行中";
    case "completed":
      return "達成！";
    case "expired":
      return "期限切れ";
  }
}

function statusColor(status: UserMission["status"]): string {
  switch (status) {
    case "active":
      return "#4CAF50";
    case "completed":
      return "#2196F3";
    case "expired":
      return "#9E9E9E";
  }
}

export default function MissionCard({ mission, onComplete, isLoading = false }: MissionCardProps) {
  const isActive = mission.status === "active";

  return (
    <div
      style={{
        border: `2px solid ${statusColor(mission.status)}`,
        borderRadius: "12px",
        padding: "16px",
        marginBottom: "12px",
        backgroundColor: mission.status === "completed" ? "#E3F2FD" : "#FFFFFF",
        opacity: mission.status === "expired" ? 0.6 : 1,
      }}
      data-testid="mission-card"
    >
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <h3 style={{ margin: 0, fontSize: "18px" }}>{mission.title}</h3>
        <span
          style={{
            backgroundColor: statusColor(mission.status),
            color: "#FFFFFF",
            padding: "4px 8px",
            borderRadius: "12px",
            fontSize: "12px",
            fontWeight: "bold",
          }}
          data-testid="mission-status"
        >
          {statusLabel(mission.status)}
        </span>
      </div>

      <p style={{ color: "#666", margin: "8px 0" }}>{mission.description}</p>

      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginTop: "12px" }}>
        <span style={{ fontSize: "14px", color: "#888" }} data-testid="mission-remaining">
          {mission.status === "completed"
            ? "観測完了"
            : formatRemainingTime(mission.valid_to)}
        </span>

        {isActive && (
          <button
            onClick={() => onComplete(mission.id)}
            disabled={isLoading}
            style={{
              backgroundColor: isLoading ? "#BDBDBD" : "#FF9800",
              color: "#FFFFFF",
              border: "none",
              borderRadius: "8px",
              padding: "8px 16px",
              fontSize: "14px",
              fontWeight: "bold",
              cursor: isLoading ? "not-allowed" : "pointer",
            }}
            data-testid="mission-complete-button"
          >
            {isLoading ? "送信中..." : "観測した！"}
          </button>
        )}
      </div>
    </div>
  );
}
