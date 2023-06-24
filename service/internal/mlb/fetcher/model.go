package fetcher

import "fmt"

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
	Status GameStatus `json:"status"`
	Teams  Teams      `json:"teams"`
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
