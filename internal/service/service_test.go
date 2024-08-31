package service

import (
	"context"
	"testing"
	"time"

	"github.com/KamilGrocholski/margo-harvester/internal/database"
)

func TestService(t *testing.T) {
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to the in-memory database: %v", err)
	}

	err = database.Migrate(db)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	service := New(db)

	startedAt := time.Now()
	endedAt := startedAt.Add(time.Hour)
	data := CreateHarvesterSessionInputData{
		"public": map[string]uint{
			"aether":  100,
			"tempest": 424,
		},
		"private": map[string]uint{
			"private_aether":  2,
			"private_tempest": 54,
		},
	}
	startedAt2 := time.Now()
	endedAt2 := startedAt.Add(time.Hour)
	data2 := CreateHarvesterSessionInputData{
		"public": map[string]uint{
			"aether":  100,
			"tempest": 424,
		},
		"private": map[string]uint{
			"private_aether":  2,
			"private_tempest": 54,
		},
	}

	ctx := context.Background()
	err = service.CreateHarvesterSession(ctx, startedAt, endedAt, data)
	if err != nil {
		t.Fatalf("Failed to create harvester session: %v", err)
	}
	err = service.CreateHarvesterSession(ctx, startedAt2, endedAt2, data2)
	if err != nil {
		t.Fatalf("Failed to create harvester session 2: %v", err)
	}

	worlds, err := service.GetAllWorlds(ctx)
	if err != nil {
		t.Fatalf("Failed to get worlds: %v", err)
	}
	if len(worlds.Worlds) != 4 {
		t.Fatalf("Expected world stats: %d, got: %d", 4, len(worlds.Worlds))
	}

	worldStatsTimeline, err := service.GetWorldStatsTimeline(
		ctx,
		"aether",
		"public",
		2,
	)
	if err != nil {
		t.Fatalf("Failed to get world timeline: %v", err)
	}

	if len(worldStatsTimeline.Timeline) != 2 {
		t.Fatalf("Expected world stats timeline timestamps: %d, got: %d", 2, len(worldStatsTimeline.Timeline))
	}

	stats := worldStatsTimeline.Timeline[1]
	if stats.Timestamp.Compare(startedAt2) != 0 {
		t.Fatalf("Expected timestamp: %v, got: %v", startedAt2, stats.Timestamp)
	}
	if stats.PlayersOnline != 100 {
		t.Fatalf("Expected players online: %v, got: %v", 100, stats.PlayersOnline)
	}
}
