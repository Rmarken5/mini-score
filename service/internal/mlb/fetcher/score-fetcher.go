package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	GameFetcher interface {
		FetchGames(time time.Time) ([]Game, error)
	}
	ScoreFetcher interface {
		FetchScore(game Game) (FetchScoreResponse, error)
	}

	Fetcher struct {
		apiURL     string
		httpClient *http.Client
	}
)

const (
	mlbAPIDomain = `https://statsapi.mlb.com`
	fetchGame    = `%s/api/v1/schedule?sportId=1,51&date=%s&gameTypes=E,S,R,A,F,D,L,W`
)

func NewFetcher(httpClient *http.Client) *Fetcher {
	return &Fetcher{httpClient: httpClient}
}

func FetchGame(fetcher GameFetcher, time time.Time) ([]Game, error) {
	return fetcher.FetchGames(time)
}

func (f *Fetcher) FetchGames(time time.Time) ([]Game, error) {
	games := make([]Game, 0)

	date := time.Format("2006-01-02")
	resp, err := f.httpClient.Get(fmt.Sprintf(fetchGame, f.apiURL, date))
	if err != nil {
		return nil, fmt.Errorf("error getting games for %s:  %w ", date, err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response:  %w ", err)
	}
	defer resp.Body.Close()

	gamesModel := &FetchGamesResponse{}
	err = json.Unmarshal(respBytes, gamesModel)
	if err != nil {
		return nil, err
	}

	if len(gamesModel.Dates) < 1 {
		return games, nil
	}

	for _, date := range gamesModel.Dates {
		games = append(games, date.Games...)
	}

	return games, nil
}

func FetchScore(fetcher ScoreFetcher, game Game) (FetchScoreResponse, error) {
	return fetcher.FetchScore(game)
}

func (f *Fetcher) FetchScore(game Game) (FetchScoreResponse, error) {
	url := fmt.Sprintf("%s%s", f.apiURL, game.Link)
	response, err := f.httpClient.Get(url)
	if err != nil {
		return FetchScoreResponse{}, err
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return FetchScoreResponse{}, err
	}
	defer response.Body.Close()

	responseScore := &FetchScoreResponse{}
	err = json.Unmarshal(responseBytes, responseScore)
	if err != nil {
		return FetchScoreResponse{}, err
	}

	return *responseScore, nil

}
