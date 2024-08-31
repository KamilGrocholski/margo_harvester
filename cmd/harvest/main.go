package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/KamilGrocholski/margo-harvester/internal/config"
	"github.com/KamilGrocholski/margo-harvester/internal/database"
	"github.com/KamilGrocholski/margo-harvester/internal/filegen"
	"github.com/KamilGrocholski/margo-harvester/internal/harvester"
	"github.com/KamilGrocholski/margo-harvester/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	err := runWithPreparation()
	if err != nil {
		os.Stdout.Write([]byte(fmt.Errorf("%v\n", err.Error()).Error()))
		os.Exit(1)
	}
}

func runWithPreparation() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	err = run(func(key string) string {
		return os.Getenv(key)
	})
	if err != nil {
		return err
	}

	return nil
}

func run(getenv config.Getenv) error {
	config, err := config.Load(getenv)
	if err != nil {
		return err
	}

	db, err := database.Open(config.Database.DB_URL)
	if err != nil {
		return err
	}
	err = database.Migrate(db)
	if err != nil {
		return err
	}

	service := service.New(db)

	httpClient := http.DefaultClient

	harvester := harvester.New(
		httpClient,
		config.HARVESTER_INTERVAL,
		config.HARVESTER_TIMEOUT,
		config.HARVESTER_MAX_ATTEMPTS,
	)

	ctx := context.Background()
	result, err := harvester.Harvest(ctx)
	if err != nil {
		return err
	}

	if areAllPlayersOnlineZero(result.Data) {
		return fmt.Errorf("all players online = 0")
	}
	err = service.CreateHarvesterSession(ctx, result.StartedAt, result.EndedAt, result.Data)
	if err != nil {
		return err
	}

	worldsList, err := service.GetAllWorlds(ctx)
	if err != nil {
		return err
	}
	err = filegen.WriteWorldsList(worldsList)
	if err != nil {
		return err
	}

	for _, world := range worldsList.Worlds {
		worldStatsTimeline, err := service.GetWorldStatsTimeline(ctx, world.Name, world.Type, 1000)
		if err != nil {
			return err
		}
		err = filegen.WriteWorldStatsTimeline(world.Name, world.Type, worldStatsTimeline)
		if err != nil {
			return err
		}
	}

	return nil
}

func areAllPlayersOnlineZero(data service.CreateHarvesterSessionInputData) bool {
	for _, stats := range data {
		for _, playersOnline := range stats {
			if playersOnline != 0 {
				return false
			}
		}
	}

	return true
}
