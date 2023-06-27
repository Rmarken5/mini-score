package fetcher

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//go:embed test-data/game.json
var gameResp []byte

//go:embed test-data/scores.json
var scoresResp []byte

func TestFetcher_FetchGames(t *testing.T) {
	testCases := map[string]struct {
		mockHttpClient   func() *httptest.Server
		expectedResponse []Game
	}{
		"should marshal games to model": {
			mockHttpClient: func() *httptest.Server {
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Define the mock response
					response := gameResp

					// Set the response status code and content type
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")

					// Write the response body
					_, err := w.Write(response)
					assert.NoError(t, err)

				}))
				return mockServer
			},
			expectedResponse: []Game{
				{
					GamePk: 717649,
					Link:   "/api/v1.1/game/717649/feed/live",
				}, {
					GamePk: 717647,
					Link:   "/api/v1.1/game/717647/feed/live",
				},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			s := tc.mockHttpClient()

			fetcher := Fetcher{s.URL, s.Client()}
			games, err := FetchGame(&fetcher, time.Now())

			assert.NoError(t, err)
			assert.EqualValues(t, tc.expectedResponse, games)
		})
	}

}

func TestFetcher_FetchScore(t *testing.T) {
	tyme, err := time.Parse("2006-01-02 3:04 PM", "2023-06-22 5:05 PM")
	require.NoError(t, err)
	testCases := map[string]struct {
		mockHttpClient   func() *httptest.Server
		expectedResponse FetchScoreResponse
	}{
		"should return linescore from the game": {
			mockHttpClient: func() *httptest.Server {
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Define the mock response
					response := scoresResp

					// Set the response status code and content type
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")

					// Write the response body
					_, err := w.Write(response)
					assert.NoError(t, err)

				}))
				return mockServer
			},
			expectedResponse: FetchScoreResponse{
				LiveData: LiveData{
					Linescore: Linescore{
						CurrentInning:        9,
						CurrentInningOrdinal: "9th",
						InningState:          "Bottom",
						InningHalf:           "Bottom",
						IsTopInning:          false,
						ScheduledInnings:     9,
						Teams: TeamStats{
							Home: TeamStat{
								Runs:       3,
								Hits:       9,
								Errors:     2,
								LeftOnBase: 6,
							},
							Away: TeamStat{
								Runs:       5,
								Hits:       9,
								Errors:     0,
								LeftOnBase: 6,
							},
						},
						Balls:   3,
						Strikes: 2,
						Outs:    3,
					}},
				GameData: GameData{
					Status: GameStatus{
						AbstractGameState: "Final",
						CodedGameState:    "F",
						DetailedState:     "Final",
						StatusCode:        "F",
						StartTimeTBD:      false,
						AbstractGameCode:  "F",
					},
					Teams: Teams{
						Away: TeamData{
							Name:          "Arizona Diamondbacks",
							Abbreviation:  "AZ",
							TeamName:      "D-backs",
							ShortName:     "Arizona",
							FranchiseName: "Arizona",
							ClubName:      "Diamondbacks",
							Active:        true,
						},
						Home: TeamData{
							Name:          "Washington Nationals",
							Abbreviation:  "WSH",
							TeamName:      "Nationals",
							ShortName:     "Washington",
							FranchiseName: "Washington",
							ClubName:      "Nationals",
							Active:        true,
						},
					},
					DateTime: DateTime{
						DateTime:     tyme,
						OriginalDate: "2023-06-22",
						OfficialDate: "2023-06-22",
						DayNight:     "day",
						Time:         "1:05",
						AMPM:         "PM",
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			s := tc.mockHttpClient()

			fetcher := Fetcher{s.URL, s.Client()}
			score, err := FetchScore(&fetcher, Game{})

			assert.NoError(t, err)
			assert.EqualValues(t, tc.expectedResponse, score)
		})
	}

}
