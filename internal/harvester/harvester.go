package harvester

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KamilGrocholski/margo-harvester/internal/service"
	"github.com/PuerkitoBio/goquery"
)

type Harvester interface {
	Harvest(context.Context) (*HarvestResult, error)
}

type harvester struct {
	httpClient         *http.Client
	requestDelay       time.Duration
	requestTimeout     time.Duration
	requestMaxAttempts int
}

func New(
	httpClient *http.Client,
	requestDelay time.Duration,
	requestTimeout time.Duration,
	requestMaxAttempts int,
) Harvester {
	return &harvester{
		httpClient:         httpClient,
		requestDelay:       requestDelay,
		requestTimeout:     requestTimeout,
		requestMaxAttempts: requestMaxAttempts,
	}
}

type HarvestResult struct {
	Data      service.CreateHarvesterSessionInputData
	StartedAt time.Time
	EndedAt   time.Time
}

func (h *harvester) Harvest(ctx context.Context) (*HarvestResult, error) {
	state := struct {
		startedAt time.Time
	}{
		startedAt: time.Now(),
	}

	response, err := h.fetchWorldsStatsPage(ctx)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	publicWorlds, err := thrashWorldsByWorldType(ctx, doc, selectorPublicWorldsList)
	if err != nil {
		return nil, err
	}
	// privateWorlds, err := thrashWorldsByWorldType(ctx, doc, selectorPrivateWorldsList)
	// if err != nil {
	// 	return nil, err
	// }

	return &HarvestResult{
		Data: service.CreateHarvesterSessionInputData{
			"public": publicWorlds,
			// "private": privateWorlds,
		},
		StartedAt: state.startedAt,
		EndedAt:   time.Now(),
	}, nil
}

func (h *harvester) fetchWorldsStatsPage(ctx context.Context) (*http.Response, error) {
	const url = "https://www.margonem.pl/stats"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not 200")
	}

	return res, nil
}

const (
	selectorPublicWorldsList   = ".public-world-list"
	selectorPrivateWorldsList  = ".private-world-list"
	selectorWorldsStats        = ".server-stats"
	selectorWorldName          = "data-name"
	selectorWorldPlayersOnline = "data-online"
)

func thrashWorldsByWorldType(
	_ context.Context, // how can i use this?
	doc *goquery.Document,
	worldsListSelector string,
) (map[string]uint, error) {
	var err error
	listEl := doc.Find(worldsListSelector)
	if listEl == nil {
		return nil, fmt.Errorf("worlds list not found")
	}
	statsEls := listEl.Find(selectorWorldsStats)
	if statsEls == nil {
		return nil, fmt.Errorf("worlds stats not found")
	}

	out := make(map[string]uint)

	statsEls.EachWithBreak(func(index int, selection *goquery.Selection) bool {
		worldName, ok := selection.Attr(selectorWorldName)
		if !ok {
			err = fmt.Errorf("world name not found")
			return false
		}
		worldName = strings.TrimPrefix(strings.TrimSpace(worldName), "#")

		playersOnline, ok := selection.Attr(selectorWorldPlayersOnline)
		if !ok {
			err = fmt.Errorf("players online not found")
			return false
		}
		parsedPlayersOnline, err := strconv.Atoi(playersOnline)
		if err != nil {
			err = fmt.Errorf("players online parsing failed")
			return false
		}

		out[worldName] = uint(parsedPlayersOnline)

		return true
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}
