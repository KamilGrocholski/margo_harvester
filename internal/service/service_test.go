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
	startedAt2 := time.Now().Add(time.Minute)
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

	stats := worldStatsTimeline.Timeline[0]
	if stats[0] != startedAt.Unix() {
		t.Fatalf("Expected timestamp 1: %v, got: %v", startedAt.Unix(), stats[0])
	}
	if stats[1] != 100 {
		t.Fatalf("Expected players online 1: %v, got: %v", 100, stats[0])
	}

	stats2 := worldStatsTimeline.Timeline[1]
	if stats2[0] != startedAt2.Unix() {
		t.Fatalf("Expected timestamp 2: %v, got: %v", startedAt2.Unix(), stats2[0])
	}
	if stats2[1] != 100 {
		t.Fatalf("Expected players online 2: %v, got: %v", 100, stats2[1])
	}
}

func TestQuerySessions(t *testing.T) {
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to the in-memory database: %v", err)
	}

	err = database.Migrate(db)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	service := New(db)

	ctx := context.Background()
	startedAt := time.Now()
	endedAt := startedAt.Add(time.Hour)
	err = service.CreateHarvesterSession(
		ctx,
		startedAt,
		endedAt,
		CreateHarvesterSessionInputData{},
	)
	if err != nil {
		t.Fatalf("Failed creating session: %v", err)
	}

	sessionsList, err := service.GetAllHarvesterSessions(ctx)
	if err != nil {
		t.Fatalf("Failed get session: %v", err)
	}

	if len(sessionsList.HarvesterSessions) != 1 {
		t.Fatalf("Expected sessions: %v, got: %v", 1, len(sessionsList.HarvesterSessions))
	}

	session := sessionsList.HarvesterSessions[0]
	if session.StartedAt.Compare(startedAt) != 0 {
		t.Fatalf("Expected startedAt: %v, got: %v", startedAt, session.StartedAt)
	}
	if session.EndedAt.Compare(endedAt) != 0 {
		t.Fatalf("Expected endedAt: %v, got: %v", endedAt, session.EndedAt)
	}
}
