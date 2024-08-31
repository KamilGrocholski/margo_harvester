package service

import (
	"context"
	"time"

	"github.com/KamilGrocholski/margo-harvester/internal/model"
	"gorm.io/gorm"
)

// WorldType -> WorldName -> PlayersOnline
type CreateHarvesterSessionInputData = map[string]map[string]uint

type HarvesterSession struct {
	ID        uint      `json:"id"`
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`
}

type HarvesterSessionsList struct {
	HarvesterSessions []HarvesterSession `json:"harvesterSessions"`
}

type WorldStatsTimeline struct {
	Timeline [][2]int64 `json:"l"` // to minify json
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

	GetAllHarvesterSessions(
		ctx context.Context,
	) (*HarvesterSessionsList, error)
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
		err := tx.
			Where("name = ?", worldType).
			FirstOrCreate(&worldTypeModel, model.WorldType{Name: worldType}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		for worldName, playersOnline := range stats {
			var worldModel model.World
			err := tx.
				Where("name = ? AND world_type_id = ?", worldName, worldTypeModel.ID).
				FirstOrCreate(&worldModel, model.World{
					Name:        worldName,
					WorldTypeID: worldTypeModel.ID,
				}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			worldStats := model.WorldStats{
				PlayersOnline:      playersOnline,
				Timestamp:          startedAt,
				WorldID:            worldModel.ID,
				HarvesterSessionID: harvesterSession.ID,
			}

			err = tx.Create(&worldStats).Error
			if err != nil {
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
		Joins("JOIN worlds ON worlds.id = world_stats.world_id").
		Joins("JOIN world_types ON world_types.id = worlds.world_type_id").
		Where("worlds.name = ? AND world_types.name = ?", worldName, worldType).
		Limit(limit).
		Find(&worldStats).
		Error
	if err != nil {
		return nil, err
	}

	timeline := make([][2]int64, len(worldStats))
	for i, s := range worldStats {
		timeline[i] = [2]int64{
			s.Timestamp.Unix(),
			int64(s.PlayersOnline),
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

func (s *service) GetAllHarvesterSessions(
	ctx context.Context,
) (*HarvesterSessionsList, error) {
	var harvesterSessions []model.HarvesterSession

	err := s.db.WithContext(ctx).
		Find(&harvesterSessions).
		Error
	if err != nil {
		return nil, err
	}

	out := make([]HarvesterSession, len(harvesterSessions))
	for i, s := range harvesterSessions {
		out[i] = HarvesterSession{
			ID:        s.ID,
			StartedAt: s.StartedAt,
			EndedAt:   s.EndedAt,
		}
	}

	return &HarvesterSessionsList{
		HarvesterSessions: out,
	}, nil
}
