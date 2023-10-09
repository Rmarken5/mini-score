package fetcher

import (
	_ "embed"
	"fmt"
	"github.com/stretchr/testify/assert"
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
					// Define the mocks response
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

	testCases := map[string]struct {
		mockHttpClient   func() *httptest.Server
		expectedResponse FetchScoreResponse
	}{
		"should return linescore from the game": {
			mockHttpClient: func() *httptest.Server {
				mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Define the mocks response
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
						Innings: []Inning{
							{
								Num:        1,
								OrdinalNum: "1st",
								Home: Home{
									Runs:       0,
									Hits:       0,
									Errors:     1,
									LeftOnBase: 0,
								},
								Away: Away{
									Runs:       1,
									Hits:       2,
									Errors:     0,
									LeftOnBase: 1,
								},
							},
							{
								Num:        2,
								OrdinalNum: "2nd",
								Home: Home{
									Runs:       0,
									Hits:       2,
									Errors:     0,
									LeftOnBase: 2,
								},
								Away: Away{
									Runs:       0,
									Hits:       0,
									Errors:     0,
									LeftOnBase: 0,
								},
							},
							{
								Num:        3,
								OrdinalNum: "3rd",
								Home: Home{
									Runs:       1,
									Hits:       1,
									Errors:     0,
									LeftOnBase: 1,
								},
								Away: Away{
									Runs:       0,
									Hits:       0,
									Errors:     0,
									LeftOnBase: 0,
								},
							},
							{
								Num:        4,
								OrdinalNum: "4th",
								Home: Home{
									Runs:       0,
									Hits:       2,
									Errors:     0,
									LeftOnBase: 2,
								},
								Away: Away{
									Runs:       1,
									Hits:       2,
									Errors:     0,
									LeftOnBase: 1,
								},
							},
							{
								Num:        5,
								OrdinalNum: "5th",
								Home: Home{
									Runs:       0,
									Hits:       1,
									Errors:     0,
									LeftOnBase: 0,
								},
								Away: Away{
									Runs:       0,
									Hits:       1,
									Errors:     0,
									LeftOnBase: 1,
								},
							},
							{
								Num:        6,
								OrdinalNum: "6th",
								Home: Home{
									Runs:       0,
									Hits:       0,
									Errors:     0,
									LeftOnBase: 0,
								},
								Away: Away{
									Runs:       0,
									Hits:       0,
									Errors:     0,
									LeftOnBase: 1,
								},
							},
							{
								Num:        7,
								OrdinalNum: "7th",
								Home: Home{
									Runs:       0,
									Hits:       1,
									Errors:     1,
									LeftOnBase: 1,
								},
								Away: Away{
									Runs:       3,
									Hits:       1,
									Errors:     0,
									LeftOnBase: 0,
								},
							},
							{
								Num:        8,
								OrdinalNum: "8th",
								Home: Home{
									Runs:       0,
									Hits:       0,
									Errors:     0,
									LeftOnBase: 0,
								},
								Away: Away{
									Runs:       0,
									Hits:       1,
									Errors:     0,
									LeftOnBase: 1,
								},
							},
							{
								Num:        9,
								OrdinalNum: "9th",
								Home: Home{
									Runs:       2,
									Hits:       2,
									Errors:     0,
									LeftOnBase: 0,
								},
								Away: Away{
									Runs:       0,
									Hits:       2,
									Errors:     0,
									LeftOnBase: 1,
								},
							},
						},
					},
				},
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
						DateTime:     time.Date(2023, 6, 22, 17, 5, 0, 0, time.UTC),
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

			fmt.Printf("%+v\n", score)
		})
	}
}
