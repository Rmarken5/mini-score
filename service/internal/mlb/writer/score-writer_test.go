package writer

import (
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	expectedResponse := fetcher.FetchScoreResponse{
		LiveData: fetcher.LiveData{Linescore: fetcher.Linescore{
			CurrentInning:        9,
			CurrentInningOrdinal: "9th",
			InningState:          "Bottom",
			InningHalf:           "Bottom",
			IsTopInning:          false,
			ScheduledInnings:     9,
			Teams: fetcher.TeamStats{
				Home: fetcher.TeamStat{
					Runs:       3,
					Hits:       9,
					Errors:     2,
					LeftOnBase: 6,
				},
				Away: fetcher.TeamStat{
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
		GameData: fetcher.GameData{
			Status: fetcher.GameStatus{
				AbstractGameState: "Final",
				CodedGameState:    "F",
				DetailedState:     "Final",
				StatusCode:        "F",
				StartTimeTBD:      false,
				AbstractGameCode:  "F",
			},
			Teams: fetcher.Teams{
				Away: fetcher.TeamData{
					Name:          "Arizona Diamondbacks",
					Abbreviation:  "AZ",
					TeamName:      "D-backs",
					ShortName:     "Arizona",
					FranchiseName: "Arizona",
					ClubName:      "Diamondbacks",
					Active:        true,
				},
				Home: fetcher.TeamData{
					Name:          "Washington Nationals",
					Abbreviation:  "WSH",
					TeamName:      "Nationals",
					ShortName:     "Washington",
					FranchiseName: "Washington",
					ClubName:      "Nationals",
					Active:        true,
				},
			},
		},
	}

	write, err := (&Writer{}).Write(&expectedResponse)
	if err != nil {
		return
	}

	expectedString :=
		`* * * * * * * * * *
*         R  H  E *
* AZ      5  9  0 *
* Final           *
* WSH     3  9  2 *
* * * * * * * * * *
`

	assert.Equal(t, expectedString, string(write))
}
