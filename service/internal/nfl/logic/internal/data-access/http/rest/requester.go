package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	scoreboardURL = "https://cdn.espn.com/core/nfl/scoreboard?xhr=1&limit=50"
)
//go:generate mockgen -destination ./schedule_requestor_mock.go -package rest . Requester
type (
	Requester interface {
		GetScoreboard() (ScoreboardResponse, error)
	}
	RequesterImpl struct {
		logger     zerolog.Logger
		httpClient *http.Client
	}
)

func NewRequester(logger zerolog.Logger, httpClient *http.Client) *RequesterImpl {
	l := logger.With().Str("service", "requester").Logger()
	return &RequesterImpl{
		logger:     l,
		httpClient: httpClient,
	}
}

func (r *RequesterImpl) GetScoreboard() (ScoreboardResponse, error) {
	logger := r.logger.With().Str("method", "GetScoreboard").Logger()
	logger.Info().Msgf("getting scoreboard")

	resp, err := r.httpClient.Get(scoreboardURL)
	if err != nil {
		logger.Error().Err(err).Msgf("while making request to: %s", scoreboardURL)
		return ScoreboardResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error().Msgf("status code %d while making request to %s", resp.StatusCode, scoreboardURL)
		return ScoreboardResponse{}, fmt.Errorf("cannot continue. response code: %d", resp.StatusCode)
	}

	var scoreboardResp ScoreboardResponse
	err = json.NewDecoder(resp.Body).Decode(&scoreboardResp)
	if err != nil {
		logger.Error().Err(err).Msg("while decoding response")
		return ScoreboardResponse{}, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("While closing response body")
		}
	}()

	return scoreboardResp, nil

}
