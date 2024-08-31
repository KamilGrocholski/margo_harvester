package database

import (
	"github.com/KamilGrocholski/margo-harvester/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(dbUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.HarvesterSession{},
		&model.World{},
		&model.WorldType{},
		&model.WorldStats{},
	)
	if err != nil {
		return err
	}

	return nil
}
