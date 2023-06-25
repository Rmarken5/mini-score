package facade

import (
	"context"
	"fmt"
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"github.com/rmarken5/mini-score/service/internal/mlb/writer"
	"sort"
	"sync"
	"time"
)

type (
	ScoreFacade interface {
		processScores(ctx context.Context, date time.Time) (string, error)
	}

	ScoreFacadeImpl struct {
		gameFetcher  fetcher.GameFetcher
		scoreFetcher fetcher.ScoreFetcher
	}
)

func NewScoreFacadeImpl(gameFetcher fetcher.GameFetcher, scoreFetcher fetcher.ScoreFetcher) *ScoreFacadeImpl {
	return &ScoreFacadeImpl{
		gameFetcher:  gameFetcher,
		scoreFetcher: scoreFetcher,
	}
}

func ProcessScores(facade ScoreFacade, ctx context.Context, date time.Time) (string, error) {
	return facade.processScores(ctx, date)
}

func (sf *ScoreFacadeImpl) processScores(ctx context.Context, date time.Time) (string, error) {
	gamesPerLine := 2
	if !IsMobile(ctx) {
		gamesPerLine = 5
	}

	games, err := sf.gameFetcher.FetchGames(date)
	if err != nil {
		return "", err
	}

	var scores []*fetcher.FetchScoreResponse
	var wg = sync.WaitGroup{}
	mutex := sync.Mutex{}
	for _, game := range games {
		wg.Add(1)
		go func(game fetcher.Game) {
			score, err := sf.scoreFetcher.FetchScore(game)
			if err != nil {
				fmt.Println(err)
			}
			mutex.Lock()
			wg.Done()
			defer mutex.Unlock()
			scores = append(scores, &score)
		}(game)
	}
	wg.Wait()
	sort.Sort(fetcher.ByGameTime(scores))
	w := writer.NewPainter(gamesPerLine)
	s, err := w.Write(scores)
	if err != nil {
		return "", err
	}

	return s, nil
}
