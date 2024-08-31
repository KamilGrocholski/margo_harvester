package service

import (
	"context"
	"time"

	"github.com/KamilGrocholski/margo-harvester/internal/model"
	"gorm.io/gorm"
)

// WorldType -> WorldName -> PlayersOnline
type CreateHarvesterSessionInputData = map[string]map[string]uint

type WorldStatsTimelineTimestamp struct {
	Timestamp     time.Time `json:"timestamp"`
	PlayersOnline uint      `json:"playersOnline"`
}

type WorldStatsTimeline struct {
	Timeline []WorldStatsTimelineTimestamp `json:"timeline"`
}

type World struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type WorldsList struct {
	Worlds []World `json:"worlds"`
}

type Service interface {
	CreateHarvesterSession(
		ctx context.Context,
		startedAt,
		endedAt time.Time,
		data CreateHarvesterSessionInputData,
	) error

	GetWorldStatsTimeline(
		ctx context.Context,
		worldName,
		worldType string,
		limit int,
	) (*WorldStatsTimeline, error)

	GetAllWorlds(
		ctx context.Context,
	) (*WorldsList, error)
}

type service struct {
	db *gorm.DB
}

func New(
	db *gorm.DB,
) Service {
	return &service{
		db: db,
	}
}

func (s *service) CreateHarvesterSession(
	ctx context.Context,
	startedAt,
	endedAt time.Time,
	data CreateHarvesterSessionInputData,
) error {
	tx := s.db.WithContext(ctx).Begin()

	var harvesterSession model.HarvesterSession
	if err := tx.
		Create(&model.HarvesterSession{
			StartedAt: startedAt,
			EndedAt:   endedAt,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for worldType, stats := range data {
		var worldTypeModel model.WorldType
		if err := tx.
			Where("name = ?", worldType).
			FirstOrCreate(&worldTypeModel, model.WorldType{Name: worldType}).Error; err != nil {
			tx.Rollback()
			return err
		}

		for worldName, playersOnline := range stats {
			var worldModel model.World
			if err := tx.
				Where("name = ? AND world_type_id = ?", worldName, worldTypeModel.ID).
				FirstOrCreate(&worldModel, model.World{
					Name:        worldName,
					WorldTypeID: worldTypeModel.ID,
				}).Error; err != nil {
				tx.Rollback()
				return err
			}

			worldStats := model.WorldStats{
				PlayersOnline:      playersOnline,
				WorldID:            worldModel.ID,
				HarvesterSessionID: harvesterSession.ID,
			}

			if err := tx.Create(&worldStats).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (s *service) GetWorldStatsTimeline(
	ctx context.Context,
	worldName,
	worldType string,
	limit int,
) (*WorldStatsTimeline, error) {
	var worldStats []model.WorldStats

	err := s.db.WithContext(ctx).
		Joins("JOIN worlds ON world_stats.world_id = worlds.id").
		Joins("JOIN world_types ON worlds.world_type_id = world_types.id").
		Where("worlds.name = ? AND world_types.name = ?", worldName, worldType).
		Preload("HarvesterSession").
		Limit(limit).
		Find(&worldStats).Error
	if err != nil {
		return nil, err
	}

	timeline := make([]WorldStatsTimelineTimestamp, len(worldStats))
	for i, s := range worldStats {
		timeline[i] = WorldStatsTimelineTimestamp{
			Timestamp:     s.HarvesterSession.StartedAt,
			PlayersOnline: s.PlayersOnline,
		}
	}

	return &WorldStatsTimeline{
		Timeline: timeline,
	}, nil
}

func (s *service) GetAllWorlds(
	ctx context.Context,
) (*WorldsList, error) {
	var worlds []model.World

	err := s.db.WithContext(ctx).
		Preload("WorldType").
		Find(&worlds).Error
	if err != nil {
		return nil, err
	}

	out := make([]World, len(worlds))
	for i, w := range worlds {
		out[i] = World{
			Name: w.Name,
			Type: w.WorldType.Name,
		}
	}

	return &WorldsList{
		Worlds: out,
	}, nil
}
