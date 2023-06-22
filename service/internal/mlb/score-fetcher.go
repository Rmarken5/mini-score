package mlb

import (
	"net/http"
	"time"
)

type (
	GameFetcher interface {
		Fetch(time time.Time) ([]string, error)
	}
	ScoreFetcher interface {
		Fetch(gameID string) (*Score, error)
	}

	Score struct {
		HomeTeam      Team
		AwayTeam      Team
		CurrentInning uint8
	}

	Team struct {
		Name               string
		Runs, Hits, Errors uint8
	}

	GameFetcherImpl struct {
		httpClient *http.Client
	}
)

func FetchGame(fetcher GameFetcher, time time.Time) ([]string, error) {
	return fetcher.Fetch(time)
}

func (g *GameFetcherImpl) Fetch(time time.Time) ([]string, error) {

}
