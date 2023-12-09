package rest

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestScores_PrintScoreboard(t *testing.T) {
	start, err := time.Parse("3:04", "1:05")
	require.NoError(t, err, "while parsing value")
	now := time.Now()
	dateString := fmt.Sprintf("%s\n", now.Format(headingTimeFormat))
	testCases := map[string]struct {
		scores         Scores
		boardsPerLine  int
		expectedString string
		err            error
	}{
		"print one game": {
			scores: Scores{{
				awayTeam: team{
					name: "PIT",
					scores: []string{
						"14", "7", "10", "7",
					},
				},
				homeTeam: team{
					name: "SF",
					scores: []string{
						"10", "10", "3", "10",
					},
				},
				quarter:   "4",
				gameClock: "05:43",
				startTime: start,
			}},
			expectedString: dateString + `* * * * * * * * * * * * * 
* Q    1  2  3  4       * 
* PIT 14  7 10  7   38  * 
* Q4             05:43  * 
* SF  10 10  3 10   33  * 
* * * * * * * * * * * * * `,
			boardsPerLine: 1,
		},
		"print one game OT": {
			scores: Scores{{
				awayTeam: team{
					name: "PIT",
					scores: []string{
						"14", "7", "10", "7", "7",
					},
				},
				homeTeam: team{
					name: "SF",
					scores: []string{
						"10", "10", "3", "10", "3",
					},
				},
				quarter:   "4",
				gameClock: "05:43",
				startTime: start,
			}},
			boardsPerLine: 2,
			expectedString: dateString + `* * * * * * * * * * * * * * 
* Q    1  2  3  4  5      * 
* PIT 14  7 10  7  7   45 * 
* Q4                05:43 * 
* SF  10 10  3 10  3   36 * 
* * * * * * * * * * * * * * `,
		},
		"print two games one line": {
			scores: Scores{{
				awayTeam: team{
					name: "PIT",
					scores: []string{
						"14", "7", "10", "7",
					},
				},
				homeTeam: team{
					name: "SF",
					scores: []string{
						"10", "10", "3", "10",
					},
				},
				quarter:   "4",
				gameClock: "05:43",
				startTime: start,
			},
				{
					awayTeam: team{
						name: "BAL",
						scores: []string{
							"14", "7", "10", "7",
						},
					},
					homeTeam: team{
						name: "IND",
						scores: []string{
							"10", "10", "3", "10",
						},
					},
					quarter:   "4",
					gameClock: "05:43",
					startTime: start,
				},
			},
			expectedString: dateString + `* * * * * * * * * * * * * * * * * * * * * * * * * * 
* Q    1  2  3  4       * * Q    1  2  3  4       * 
* PIT 14  7 10  7   38  * * BAL 14  7 10  7   38  * 
* Q4             05:43  * * Q4             05:43  * 
* SF  10 10  3 10   33  * * IND 10 10  3 10   33  * 
* * * * * * * * * * * * * * * * * * * * * * * * * * `,
			boardsPerLine: 2,
		},
		"print two games two lines": {
			scores: Scores{{
				awayTeam: team{
					name: "PIT",
					scores: []string{
						"14", "7", "10", "7",
					},
				},
				homeTeam: team{
					name: "SF",
					scores: []string{
						"10", "10", "3", "10",
					},
				},
				quarter:   "4",
				gameClock: "05:43",
				startTime: start,
			},
				{
					awayTeam: team{
						name: "BAL",
						scores: []string{
							"14", "7", "10", "7",
						},
					},
					homeTeam: team{
						name: "IND",
						scores: []string{
							"10", "10", "3", "10",
						},
					},
					quarter:   "4",
					gameClock: "05:43",
					startTime: start,
				},
			},
			expectedString: dateString + `* * * * * * * * * * * * * 
* Q    1  2  3  4       * 
* PIT 14  7 10  7   38  * 
* Q4             05:43  * 
* SF  10 10  3 10   33  * 
* * * * * * * * * * * * * 
* * * * * * * * * * * * * 
* Q    1  2  3  4       * 
* BAL 14  7 10  7   38  * 
* Q4             05:43  * 
* IND 10 10  3 10   33  * 
* * * * * * * * * * * * * `,
			boardsPerLine: 1,
		},
		"print three games two lines": {
			scores: Scores{{
				awayTeam: team{
					name: "PIT",
					scores: []string{
						"14", "7", "10", "7",
					},
				},
				homeTeam: team{
					name: "SF",
					scores: []string{
						"10", "10", "3", "10",
					},
				},
				quarter:   "4",
				gameClock: "05:43",
				startTime: start,
			},
				{
					awayTeam: team{
						name: "BAL",
						scores: []string{
							"14", "7", "10", "7",
						},
					},
					homeTeam: team{
						name: "IND",
						scores: []string{
							"10", "10", "3", "10",
						},
					},
					quarter:   "4",
					gameClock: "05:43",
					startTime: start,
				},
				{
					awayTeam: team{
						name: "ATL",
						scores: []string{
							"14", "7", "10", "7",
						},
					},
					homeTeam: team{
						name: "KC",
						scores: []string{
							"10", "10", "3", "10",
						},
					},
					quarter:   "4",
					gameClock: "05:43",
					startTime: start,
				},
			},
			expectedString: dateString + `* * * * * * * * * * * * * * * * * * * * * * * * * * 
* Q    1  2  3  4       * * Q    1  2  3  4       * 
* PIT 14  7 10  7   38  * * BAL 14  7 10  7   38  * 
* Q4             05:43  * * Q4             05:43  * 
* SF  10 10  3 10   33  * * IND 10 10  3 10   33  * 
* * * * * * * * * * * * * * * * * * * * * * * * * * 
* * * * * * * * * * * * * 
* Q    1  2  3  4       * 
* ATL 14  7 10  7   38  * 
* Q4             05:43  * 
* KC  10 10  3 10   33  * 
* * * * * * * * * * * * * `,
			boardsPerLine: 2,
		},
	}
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			bw := bytes.Buffer{}
			err := tc.scores.PrintScoreboard(&bw, now, tc.boardsPerLine)
			assert.ErrorIs(t, err, tc.err)
			gotString := bw.String()
			fmt.Println(gotString)
			assert.Equal(t, tc.expectedString, gotString)
		})
	}
}
