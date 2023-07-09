package fetcher

import (
	"fmt"
	"time"
)

type FetchGamesResponse struct {
	Dates []Date `json:"dates"`
}

type Date struct {
	Games []Game `json:"games"`
}

type Game struct {
	GamePk int    `json:"gamePk"`
	Link   string `json:"link"`
}

type GameResult struct {
}

type FetchScoreResponse struct {
	LiveData LiveData `json:"liveData"`
	GameData GameData `json:"gameData"`
}

type ByGameTime []*FetchScoreResponse

func (a ByGameTime) Len() int      { return len(a) }
func (a ByGameTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByGameTime) Less(i, j int) bool {
	if a[i].GameData.DateTime.DateTime.Equal(a[j].GameData.DateTime.DateTime) {
		return a[i].GameData.Teams.Home.TeamName < a[j].GameData.Teams.Home.TeamName
	}
	return a[i].GameData.DateTime.DateTime.Before(a[j].GameData.DateTime.DateTime)
}

func (fsr *FetchScoreResponse) String() string {
	inning := ""
	away := fmt.Sprintf("* %s    %s *", fsr.GameData.Teams.Away.String(), fsr.LiveData.Linescore.Teams.Away.String())
	if fsr.GameData.Status.StatusCode == "F" {
		inning = fsr.GameData.Status.DetailedState
	} else {
		inning = fsr.LiveData.Linescore.InningHalf + " " + fsr.LiveData.Linescore.CurrentInningOrdinal
	}
	inningString := fmt.Sprintf("* %-15s *", inning)
	home := fmt.Sprintf("* %s    %s *", fsr.GameData.Teams.Home.String(), fsr.LiveData.Linescore.Teams.Home.String())

	return fmt.Sprintf("%s\n%s\n%s", away, inningString, home)
}

type LiveData struct {
	Linescore Linescore `json:"linescore"`
}

type GameData struct {
	Status   GameStatus `json:"status"`
	Teams    Teams      `json:"teams"`
	DateTime DateTime   `json:"datetime"`
}

type DateTime struct {
	DateTime     time.Time `json:"dateTime"`
	OriginalDate string    `json:"originalDate"`
	OfficialDate string    `json:"officialDate"`
	DayNight     string    `json:"dayNight"`
	Time         string    `json:"time"`
	AMPM         string    `json:"ampm"`
}

type Teams struct {
	Away TeamData `json:"away"`
	Home TeamData `json:"home"`
}
type GameStatus struct {
	AbstractGameState string `json:"abstractGameState"`
	CodedGameState    string `json:"codedGameState"`
	DetailedState     string `json:"detailedState"`
	StatusCode        string `json:"statusCode"`
	StartTimeTBD      bool   `json:"startTimeTBD"`
	AbstractGameCode  string `json:"abstractGameCode"`
}

type TeamData struct {
	Name          string `json:"name"`
	Abbreviation  string `json:"abbreviation"`
	TeamName      string `json:"teamName"`
	ShortName     string `json:"shortName"`
	FranchiseName string `json:"franchiseName"`
	ClubName      string `json:"clubName"`
	Active        bool   `json:"active"`
}

func (td *TeamData) String() string {
	return fmt.Sprintf("%-3s", td.Abbreviation)
}

type Linescore struct {
	CurrentInning        int       `json:"currentInning"`
	CurrentInningOrdinal string    `json:"currentInningOrdinal"`
	InningState          string    `json:"inningState"`
	InningHalf           string    `json:"inningHalf"`
	IsTopInning          bool      `json:"isTopInning"`
	ScheduledInnings     int       `json:"scheduledInnings"`
	Teams                TeamStats `json:"teams"`
	Balls                int       `json:"balls"`
	Strikes              int       `json:"strikes"`
	Outs                 int       `json:"outs"`
	Innings              Innings   `json:"innings"`
}
type Innings []Inning
type Inning struct {
	Num        int    `json:"num,omitempty"`
	OrdinalNum string `json:"ordinalNum,omitempty"`
	Home       Home   `json:"home,omitempty"`
	Away       Away   `json:"away,omitempty"`
}

func (i Innings) PrintInningRuns() (away string, home string) {
	iLen := len(i)
	for _, inn := range i {
		away += fmt.Sprintf(" %2d", inn.Away.Runs)
		home += fmt.Sprintf(" %2d", inn.Home.Runs)
	}
	if iLen < 9 {
		for k := 0; k < 9-iLen; k++ {
			away += "   "
			home += "   "
		}
	}
	return // Naked return
}

type Home struct {
	Runs       int `json:"runs,omitempty"`
	Hits       int `json:"hits,omitempty"`
	Errors     int `json:"errors,omitempty"`
	LeftOnBase int `json:"leftOnBase,omitempty"`
}
type Away struct {
	Runs       int `json:"runs,omitempty"`
	Hits       int `json:"hits,omitempty"`
	Errors     int `json:"errors,omitempty"`
	LeftOnBase int `json:"leftOnBase,omitempty"`
}

type TeamStat struct {
	Runs       int `json:"runs"`
	Hits       int `json:"hits"`
	Errors     int `json:"errors"`
	LeftOnBase int `json:"leftOnBase"`
}

type TeamStats struct {
	Home TeamStat `json:"home"`
	Away TeamStat `json:"away"`
}

func (ts *TeamStat) String() string {
	return fmt.Sprintf("%2d %2d %2d", ts.Runs, ts.Hits, ts.Errors)
}
