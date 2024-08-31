package service

import (
	"context"
	"time"
)

type MockService struct {
	CreateHarvesterSessionFunc func(
		ctx context.Context,
		startedAt,
		endedAt time.Time,
		data CreateHarvesterSessionInputData,
	) error

	GetWorldStatsTimelineFunc func(
		ctx context.Context,
		worldName,
		worldType string,
		limit int,
	) (*WorldStatsTimeline, error)

	GetAllWorldsFunc func(
		ctx context.Context,
	) (*WorldsList, error)
}

func (m MockService) CreateHarvesterSession(
	ctx context.Context,
	startedAt,
	endedAt time.Time,
	data CreateHarvesterSessionInputData,
) error {
	return m.CreateHarvesterSessionFunc(ctx, startedAt, endedAt, data)
}

func (m MockService) GetWorldStatsTimeline(
	ctx context.Context,
	worldName,
	worldType string,
	limit int,
) (*WorldStatsTimeline, error) {
	return m.GetWorldStatsTimelineFunc(ctx, worldName, worldType, limit)
}

func (m MockService) GetAllWorlds(ctx context.Context) (*WorldsList, error) {
	return m.GetAllWorldsFunc(ctx)
}
