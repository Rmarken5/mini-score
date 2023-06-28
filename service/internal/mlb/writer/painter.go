package writer

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"strings"
	"text/template"
	"time"
)

//go:embed templates/game-time.gotmpl
var gtTemplate embed.FS

var layout = "Jan, 02 2006"

type (
	Painter struct {
		date             time.Time
		gamesPerLine     int
		Games            int
		AwayTeamLine     []string
		GameProgressLine []string
		HomeTeamLine     []string
	}
)

func NewPainter(gamesPerLine int, date time.Time) *Painter {
	return &Painter{gamesPerLine: gamesPerLine, date: date}
}

func (p *Painter) addScore(score *fetcher.FetchScoreResponse) {
	p.Games++
	p.AwayTeamLine = append(p.AwayTeamLine, fmt.Sprintf(" * %s    %s * ", score.GameData.Teams.Away.String(), score.LiveData.Linescore.Teams.Away.String()))
	p.HomeTeamLine = append(p.HomeTeamLine, fmt.Sprintf(" * %s    %s * ", score.GameData.Teams.Home.String(), score.LiveData.Linescore.Teams.Home.String()))
	var gameStatus string
	switch score.GameData.Status.StatusCode {
	case "F":
		gameStatus = score.GameData.Status.DetailedState
	case "P", "S":
		gameStatus = fmt.Sprintf("%s %s", score.GameData.DateTime.Time, score.GameData.DateTime.AMPM)
	default:
		gameStatus = fmt.Sprintf("%s %s", score.LiveData.Linescore.InningHalf, score.LiveData.Linescore.CurrentInningOrdinal)
	}

	gameStatusString := fmt.Sprintf(" * %-15s * ", gameStatus)

	p.GameProgressLine = append(p.GameProgressLine, gameStatusString)

}

func (p *Painter) Write(scores []*fetcher.FetchScoreResponse) (string, error) {
	for _, score := range scores {
		p.addScore(score)
	}

	sb := strings.Builder{}

	for i := 0; i < p.Games; i += p.gamesPerLine {
		var limit int = p.gamesPerLine
		if p.Games < i+p.gamesPerLine {
			limit = p.gamesPerLine - ((i + p.gamesPerLine) - p.Games)
		}

		for j := 0; j < limit; j++ {
			sb.WriteString(" * * * * * * * * * * ")
		}
		sb.Write([]byte("\n"))
		for j := 0; j < limit; j++ {
			sb.WriteString(" *         R  H  E * ")
		}
		sb.Write([]byte("\n"))
		for j := 0; j < limit; j++ {
			sb.WriteString(p.AwayTeamLine[i+j])
		}
		sb.Write([]byte("\n"))
		for j := 0; j < limit; j++ {
			sb.WriteString(p.GameProgressLine[i+j])
		}
		sb.Write([]byte("\n"))
		for j := 0; j < limit; j++ {
			sb.WriteString(p.HomeTeamLine[i+j])
		}
		sb.Write([]byte("\n"))
		for j := 0; j < limit; j++ {
			sb.WriteString(" * * * * * * * * * * ")
		}
		sb.Write([]byte("\n"))
	}

	file, err := template.ParseFS(gtTemplate, "templates/game-time.gotmpl")
	if err != nil {
		return "", err
	}

	buff := bytes.NewBuffer(nil)

	err = file.Execute(buff, struct {
		Time  string
		Games string
	}{
		Time:  p.date.Format(layout),
		Games: sb.String(),
	})
	if err != nil {
		return "", err
	}

	return buff.String(), nil

}