package facade

import (
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"github.com/rmarken5/mini-score/service/internal/mlb/writer"
	"time"
)

type (
	ScoreFacade interface {
		ProcessScores(date time.Time) (string, error)
	}

	ScoreFacadeImpl struct {
		gameFetcher  fetcher.GameFetcher
		scoreFetcher fetcher.ScoreFetcher
		writer       writer.Write
	}
)

func NewScoreFacadeImpl(gameFetcher fetcher.GameFetcher, scoreFetcher fetcher.ScoreFetcher, write writer.Write) *ScoreFacadeImpl {
	return &ScoreFacadeImpl{
		gameFetcher:  gameFetcher,
		scoreFetcher: scoreFetcher,
		writer:       write,
	}
}

func (sf *ScoreFacadeImpl) ProcessScores(date time.Time) (string, error) {
	games, err := sf.gameFetcher.FetchGames(date)
	if err != nil {
		return "", err
	}

	var scores []*fetcher.FetchScoreResponse
	for _, game := range games {
		score, err := sf.scoreFetcher.FetchScore(game)
		if err != nil {
			return "", err
		}
		scores = append(scores, &score)
	}

	s, err := sf.writer.Write(scores)
	if err != nil {
		return "", err
	}

	return s, nil
}
