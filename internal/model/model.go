package model

import (
	"time"

	"gorm.io/gorm"
)

type HarvesterSession struct {
	gorm.Model
	StartedAt time.Time
	EndedAt   time.Time

	WorldStats []WorldStats `gorm:"foreignKey:HarvesterSessionID"`
}

type World struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;not null"`

	WorldTypeID uint `gorm:"index"`
	WorldType   WorldType
	WorldsStats []WorldStats `gorm:"foreignKey:WorldID"`
}

type WorldType struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
}

type WorldStats struct {
	gorm.Model
	PlayersOnline uint

	WorldID            uint `gorm:"index"`
	World              World
	HarvesterSessionID uint `gorm:"index"`
	HarvesterSession   HarvesterSession
}
