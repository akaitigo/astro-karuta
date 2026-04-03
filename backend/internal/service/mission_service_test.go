package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/repository"
	"github.com/akaitigo/astro-karuta/backend/internal/seed"
	"github.com/akaitigo/astro-karuta/backend/internal/service"
)

func setupMissionService(t *testing.T) (*service.MissionService, *repository.InMemoryCardRepository) {
	t.Helper()
	cardRepo := repository.NewInMemoryCardRepository()
	missionRepo := repository.NewInMemoryMissionRepository()

	if err := seed.LoadCards(context.Background(), cardRepo); err != nil {
		t.Fatal(err)
	}
	svc := service.NewMissionService(missionRepo, cardRepo)
	return svc, cardRepo
}

func TestMissionService_GetActiveMissions_EmptyUserID(t *testing.T) {
	svc, _ := setupMissionService(t)

	_, err := svc.GetActiveMissions(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty user ID")
	}
}

func TestMissionService_GetActiveMissions_GeneratesMissions(t *testing.T) {
	svc, _ := setupMissionService(t)

	// Set time to July (summer) for predictable results
	svc.SetNowFunc(func() time.Time {
		return time.Date(2026, 7, 15, 20, 0, 0, 0, time.Local)
	})

	missions, err := svc.GetActiveMissions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(missions) == 0 {
		t.Fatal("expected missions to be generated")
	}
	if len(missions) > 5 {
		t.Errorf("expected at most 5 missions, got %d", len(missions))
	}

	for _, m := range missions {
		if m.UserID != "user-1" {
			t.Errorf("expected user ID user-1, got %s", m.UserID)
		}
		if m.Status != "active" {
			t.Errorf("expected active status, got %s", m.Status)
		}
	}
}

func TestMissionService_GetActiveMissions_ReturnsCached(t *testing.T) {
	svc, _ := setupMissionService(t)

	svc.SetNowFunc(func() time.Time {
		return time.Date(2026, 7, 15, 20, 0, 0, 0, time.Local)
	})

	// First call generates
	missions1, err := svc.GetActiveMissions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Second call returns same missions
	missions2, err := svc.GetActiveMissions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(missions1) != len(missions2) {
		t.Errorf("expected same count: %d vs %d", len(missions1), len(missions2))
	}
}

func TestMissionService_CompleteMission_Success(t *testing.T) {
	svc, _ := setupMissionService(t)

	// Generate missions at nighttime in July
	svc.SetNowFunc(func() time.Time {
		return time.Date(2026, 7, 15, 21, 0, 0, 0, time.Local)
	})

	missions, err := svc.GetActiveMissions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(missions) == 0 {
		t.Fatal("expected missions to be generated")
	}

	// Complete the first mission with valid location (near Tokyo)
	resp, err := svc.CompleteMission(context.Background(), "user-1", missions[0].ID, 35.0, 139.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Mission.Status != "completed" {
		t.Errorf("expected completed status, got %s", resp.Mission.Status)
	}
	if resp.BonusCard == nil {
		t.Error("expected bonus card")
	}
	if resp.Mission.CompletedAt == nil {
		t.Error("expected completed_at to be set")
	}
}

func TestMissionService_CompleteMission_EmptyUserID(t *testing.T) {
	svc, _ := setupMissionService(t)

	_, err := svc.CompleteMission(context.Background(), "", "mission-1", 35.0, 139.0)
	if err == nil {
		t.Fatal("expected error for empty user ID")
	}
}

func TestMissionService_CompleteMission_EmptyMissionID(t *testing.T) {
	svc, _ := setupMissionService(t)

	_, err := svc.CompleteMission(context.Background(), "user-1", "", 35.0, 139.0)
	if err == nil {
		t.Fatal("expected error for empty mission ID")
	}
}

func TestMissionService_CompleteMission_WrongUser(t *testing.T) {
	svc, _ := setupMissionService(t)

	svc.SetNowFunc(func() time.Time {
		return time.Date(2026, 7, 15, 21, 0, 0, 0, time.Local)
	})

	missions, err := svc.GetActiveMissions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Try to complete as a different user
	_, err = svc.CompleteMission(context.Background(), "user-2", missions[0].ID, 35.0, 139.0)
	if err == nil {
		t.Fatal("expected error for wrong user")
	}
}

func TestMissionService_CompleteMission_Daytime(t *testing.T) {
	svc, _ := setupMissionService(t)

	// Generate at night
	svc.SetNowFunc(func() time.Time {
		return time.Date(2026, 7, 15, 21, 0, 0, 0, time.Local)
	})

	missions, err := svc.GetActiveMissions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Try to complete during daytime
	svc.SetNowFunc(func() time.Time {
		return time.Date(2026, 7, 15, 12, 0, 0, 0, time.Local)
	})

	_, err = svc.CompleteMission(context.Background(), "user-1", missions[0].ID, 35.0, 139.0)
	if err == nil {
		t.Fatal("expected error for daytime completion")
	}
}

func TestMissionService_CompleteMission_LocationTooFar(t *testing.T) {
	svc, _ := setupMissionService(t)

	svc.SetNowFunc(func() time.Time {
		return time.Date(2026, 7, 15, 21, 0, 0, 0, time.Local)
	})

	missions, err := svc.GetActiveMissions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Location far from Tokyo (latitude deviation > 30)
	_, err = svc.CompleteMission(context.Background(), "user-1", missions[0].ID, -20.0, 139.0)
	if err == nil {
		t.Fatal("expected error for location too far")
	}
}

func TestMissionService_CompleteMission_AlreadyCompleted(t *testing.T) {
	svc, _ := setupMissionService(t)

	svc.SetNowFunc(func() time.Time {
		return time.Date(2026, 7, 15, 21, 0, 0, 0, time.Local)
	})

	missions, err := svc.GetActiveMissions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Complete once
	_, err = svc.CompleteMission(context.Background(), "user-1", missions[0].ID, 35.0, 139.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Try to complete again
	_, err = svc.CompleteMission(context.Background(), "user-1", missions[0].ID, 35.0, 139.0)
	if err == nil {
		t.Fatal("expected error for already completed mission")
	}
}

func TestMissionService_CompleteMission_NighttimeBoundary(t *testing.T) {
	svc, _ := setupMissionService(t)

	// Generate missions
	svc.SetNowFunc(func() time.Time {
		return time.Date(2026, 7, 15, 22, 0, 0, 0, time.Local)
	})

	missions, err := svc.GetActiveMissions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	missionID := missions[0].ID

	testCases := []struct {
		name    string
		hour    int
		wantErr bool
	}{
		{"midnight", 0, false},
		{"5am", 5, false},
		{"6am_rejected", 6, true},
		{"noon", 12, true},
		{"5pm_rejected", 17, true},
		{"6pm_accepted", 18, false},
		{"11pm", 23, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Re-setup to get fresh missions
			innerSvc, _ := setupMissionService(t)
			innerSvc.SetNowFunc(func() time.Time {
				return time.Date(2026, 7, 15, 22, 0, 0, 0, time.Local)
			})

			innerMissions, err := innerSvc.GetActiveMissions(context.Background(), "user-1")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_ = missionID // referenced from outer scope only as pattern

			innerSvc.SetNowFunc(func() time.Time {
				return time.Date(2026, 7, 15, tc.hour, 0, 0, 0, time.Local)
			})

			_, err = innerSvc.CompleteMission(context.Background(), "user-1", innerMissions[0].ID, 35.0, 139.0)
			if tc.wantErr && err == nil {
				t.Errorf("hour %d: expected error", tc.hour)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("hour %d: unexpected error: %v", tc.hour, err)
			}
		})
	}
}
