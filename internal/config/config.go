package config

import (
	"fmt"
	"strconv"
	"time"
)

type Getenv func(string) string

type Config struct {
	Database
	Harvester
}

type Database struct {
	DB_URL string
}

type Harvester struct {
	HARVESTER_TIMEOUT      time.Duration
	HARVESTER_INTERVAL     time.Duration
	HARVESTER_MAX_ATTEMPTS int
}

func Load(getenv Getenv) (Config, error) {
	config := Config{}

	database, err := LoadDatabase(getenv)
	if err != nil {
		return config, err
	}
	config.Database = database

	harvester, err := LoadHarvester(getenv)
	if err != nil {
		return config, err
	}
	config.Harvester = harvester

	return config, nil
}

func LoadDatabase(getenv Getenv) (Database, error) {
	config := Database{}

	config.DB_URL = getenv("DB_URL")
	if config.DB_URL == "" {
		return config, fmt.Errorf("no DB_URL")
	}

	return config, nil
}

func LoadHarvester(getenv Getenv) (Harvester, error) {
	config := Harvester{}

	timeout, err := time.ParseDuration(getenv("HARVESTER_TIMEOUT"))
	if err != nil {
		return config, fmt.Errorf("err reading HARVESTER_TIMEOUT: %v", err)
	}
	config.HARVESTER_TIMEOUT = timeout

	interval, err := time.ParseDuration(getenv("HARVESTER_INTERVAL"))
	if err != nil {
		return config, fmt.Errorf("err reading HARVESTER_INTERVAL: %v", err)
	}
	config.HARVESTER_INTERVAL = interval

	maxAttempts, err := strconv.Atoi(getenv("HARVESTER_MAX_ATTEMPTS"))
	if err != nil {
		return config, fmt.Errorf("err reading HARVESTER_MAX_ATTEMPTS: %v", err)
	}
	config.HARVESTER_MAX_ATTEMPTS = maxAttempts

	return config, nil
}
