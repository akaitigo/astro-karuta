"use client";

import { useEffect, useState, useCallback } from "react";
import Link from "next/link";
import MissionCard from "@/components/MissionCard";
import type { UserMission, CompleteMissionResponse } from "@/types/mission";

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

// Temporary user ID until auth is implemented
const DEFAULT_USER_ID = "00000000-0000-4000-8000-000000000001";

export default function MissionsPage() {
  const [missions, setMissions] = useState<UserMission[]>([]);
  const [loading, setLoading] = useState(true);
  const [completingId, setCompletingId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const fetchMissions = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const res = await fetch(
        `${API_BASE}/api/v1/missions?user_id=${DEFAULT_USER_ID}`
      );
      if (!res.ok) {
        throw new Error("ミッションの取得に失敗しました");
      }
      const data: UserMission[] = await res.json();
      setMissions(data ?? []);
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("不明なエラーが発生しました");
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void fetchMissions();
  }, [fetchMissions]);

  const handleComplete = useCallback(
    async (missionId: string) => {
      try {
        setCompletingId(missionId);
        setError(null);
        setSuccessMessage(null);

        // Get current position via Geolocation API
        let lat = 35.68;
        let lng = 139.77;

        if (typeof navigator !== "undefined" && navigator.geolocation) {
          try {
            const position = await new Promise<GeolocationPosition>(
              (resolve, reject) => {
                navigator.geolocation.getCurrentPosition(resolve, reject, {
                  timeout: 5000,
                });
              }
            );
            lat = position.coords.latitude;
            lng = position.coords.longitude;
          } catch {
            // Use default coordinates if geolocation fails
          }
        }

        const res = await fetch(
          `${API_BASE}/api/v1/missions/${missionId}/complete`,
          {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              user_id: DEFAULT_USER_ID,
              lat,
              lng,
            }),
          }
        );

        if (!res.ok) {
          const errData: { error?: string } = await res.json();
          throw new Error(errData.error ?? "ミッションの完了に失敗しました");
        }

        const data: CompleteMissionResponse = await res.json();

        // Update mission in local state
        setMissions((prev) =>
          prev.map((m) => (m.id === missionId ? data.mission : m))
        );

        if (data.bonus_card) {
          setSuccessMessage(
            `${data.bonus_card.name}のカードを獲得しました！`
          );
        } else {
          setSuccessMessage("ミッション達成！");
        }
      } catch (err) {
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError("不明なエラーが発生しました");
        }
      } finally {
        setCompletingId(null);
      }
    },
    []
  );

  return (
    <main style={{ padding: "24px", maxWidth: "600px", margin: "0 auto" }}>
      <h1>観測ミッション</h1>
      <p style={{ color: "#666", marginBottom: "24px" }}>
        今月見える星座を実際に観測してカードを集めよう！
      </p>

      {error && (
        <div
          style={{
            backgroundColor: "#FFEBEE",
            color: "#D32F2F",
            padding: "12px",
            borderRadius: "8px",
            marginBottom: "16px",
          }}
          data-testid="mission-error"
        >
          {error}
        </div>
      )}

      {successMessage && (
        <div
          style={{
            backgroundColor: "#E8F5E9",
            color: "#2E7D32",
            padding: "12px",
            borderRadius: "8px",
            marginBottom: "16px",
          }}
          data-testid="mission-success"
        >
          {successMessage}
        </div>
      )}

      {loading ? (
        <p>読み込み中...</p>
      ) : missions.length === 0 ? (
        <p>現在利用可能なミッションはありません。</p>
      ) : (
        <div>
          {missions.map((mission) => (
            <MissionCard
              key={mission.id}
              mission={mission}
              onComplete={(id) => void handleComplete(id)}
              isLoading={completingId === mission.id}
            />
          ))}
        </div>
      )}

      <div style={{ marginTop: "24px" }}>
        <Link
          href="/"
          style={{
            color: "#1565C0",
            textDecoration: "none",
            fontSize: "14px",
          }}
        >
          ホームに戻る
        </Link>
      </div>
    </main>
  );
}
