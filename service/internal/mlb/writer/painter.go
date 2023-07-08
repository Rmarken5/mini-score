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
		lineLength       int
		Games            int
		TopBottomBorder  []string
		InningsLine      []string
		AwayTeamLine     []string
		GameProgressLine []string
		HomeTeamLine     []string
	}
)

func NewPainter(lineLength int, date time.Time) *Painter {
	return &Painter{lineLength: lineLength, date: date}
}

func (p *Painter) addScore(score *fetcher.FetchScoreResponse) {
	p.Games++

	awayInnings, homeInnings := score.LiveData.Linescore.Innings.PrintInningRuns()

	inningLen := 9
	if inningLen < len(score.LiveData.Linescore.Innings) {
		inningLen = len(score.LiveData.Linescore.Innings)
	}
	inningHeader := " *         "
	for i := 0; i < inningLen; i++ {
		inningHeader += fmt.Sprintf("%d  ", i+1)
	}
	inningHeader += "  R  H  E  * "
	inningHeaderLen := len(inningHeader)

	// add one to make even two characters are written at a time.
	if inningHeaderLen%2 > 0 {
		inningHeaderLen++
	}

	topAndBottomBorder := ""

	for i := 0; i < (inningHeaderLen-2)/2; i++ {
		topAndBottomBorder += " *"
	}
	topAndBottomBorder += " "

	p.TopBottomBorder = append(p.TopBottomBorder, topAndBottomBorder)
	p.InningsLine = append(p.InningsLine, inningHeader)

	p.AwayTeamLine = append(p.AwayTeamLine, fmt.Sprintf(" * %s   %s   %s  * ", score.GameData.Teams.Away.String(), awayInnings, score.LiveData.Linescore.Teams.Away.String()))
	p.HomeTeamLine = append(p.HomeTeamLine, fmt.Sprintf(" * %s   %s   %s  * ", score.GameData.Teams.Home.String(), homeInnings, score.LiveData.Linescore.Teams.Home.String()))
	var gameStatus string
	switch score.GameData.Status.StatusCode {
	case "F":
		gameStatus = score.GameData.Status.DetailedState
	case "P", "S":
		gameStatus = fmt.Sprintf("%s %s", score.GameData.DateTime.Time, score.GameData.DateTime.AMPM)
	default:
		gameStatus = fmt.Sprintf("%s %s", score.LiveData.Linescore.InningHalf, score.LiveData.Linescore.CurrentInningOrdinal)
	}
	formatter := " * %-" + fmt.Sprintf("%d", len(inningHeader)-6) + "s * "
	gameStatusString := fmt.Sprintf(formatter, gameStatus)

	p.GameProgressLine = append(p.GameProgressLine, gameStatusString)

}

func (p *Painter) Write(scores []*fetcher.FetchScoreResponse) (string, error) {
	for _, score := range scores {
		p.addScore(score)
	}

	sb := strings.Builder{}

	for i := 0; i < p.Games; i += p.lineLength {
		var limit int = p.lineLength
		if p.Games < i+p.lineLength {
			limit = p.lineLength - ((i + p.lineLength) - p.Games)
		}

		for j := 0; j < limit; j++ {
			sb.WriteString(p.TopBottomBorder[i+j])
		}
		sb.Write([]byte("\n"))
		for j := 0; j < limit; j++ {
			sb.WriteString(p.InningsLine[i+j])
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
			sb.WriteString(p.TopBottomBorder[i+j])
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
