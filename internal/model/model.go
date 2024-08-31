package model

import (
	"time"

	"gorm.io/gorm"
)

type HarvesterSession struct {
	gorm.Model
	StartedAt time.Time `gorm:"not null"`
	EndedAt   time.Time `gorm:"not null"`

	WorldStats []WorldStats `gorm:"foreignKey:HarvesterSessionID"`
}

type World struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;not null"`

	WorldTypeID uint         `gorm:"index"`
	WorldType   WorldType    `gorm:"foreignKey:WorldTypeID"`
	WorldsStats []WorldStats `gorm:"foreignKey:WorldID"`
}

type WorldType struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
}

type WorldStats struct {
	gorm.Model
	PlayersOnline uint
	Timestamp     time.Time

	WorldID            uint             `gorm:"index"`
	World              World            `gorm:"foreignKey:WorldID"`
	HarvesterSessionID uint             `gorm:"index"`
	HarvesterSession   HarvesterSession `gorm:"foreignKey:WorldID"`
}
