package filegen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KamilGrocholski/margo-harvester/internal/service"
)

func GetPath(relativePath string) (string, error) {
	return filepath.Join("public", relativePath), nil
}

func WriteWorldStatsTimeline(worldName, worldType string, timeline *service.WorldStatsTimeline) error {
	dir := filepath.Join("public", "worlds", worldType)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	filePath := filepath.Join(dir, fmt.Sprintf("%s.json", worldName))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(timeline); err != nil {
		return fmt.Errorf("failed to encode timeline to JSON: %v", err)
	}

	return nil
}

func WriteWorldsList(worldsList *service.WorldsList) error {
	filePath := filepath.Join("public", "worlds.json")

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(worldsList); err != nil {
		return fmt.Errorf("failed to encode worlds list to JSON: %v", err)
	}

	return nil
}
