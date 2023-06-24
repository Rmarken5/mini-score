package writer

import (
	"fmt"
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"strings"
)

type (
	Painter struct {
		gamesPerLine     int
		Games            int
		AwayTeamLine     []string
		GameProgressLine []string
		HomeTeamLine     []string
	}
)

func NewPainter(gamesPerLine int) *Painter {
	return &Painter{gamesPerLine: gamesPerLine}
}

func (p *Painter) addScore(score *fetcher.FetchScoreResponse) {
	p.Games++
	p.AwayTeamLine = append(p.AwayTeamLine, fmt.Sprintf(" * %s    %s * ", score.GameData.Teams.Away.String(), score.LiveData.Linescore.Teams.Away.String()))
	p.HomeTeamLine = append(p.HomeTeamLine, fmt.Sprintf(" * %s    %s * ", score.GameData.Teams.Home.String(), score.LiveData.Linescore.Teams.Home.String()))
	var inning string
	if score.GameData.Status.StatusCode == "F" {
		inning = score.GameData.Status.DetailedState
	} else {
		inning = score.LiveData.Linescore.InningHalf + " " + score.LiveData.Linescore.CurrentInningOrdinal
	}
	inningString := fmt.Sprintf(" * %-15s * ", inning)

	p.GameProgressLine = append(p.GameProgressLine, inningString)

}

func (p *Painter) Write(scores []*fetcher.FetchScoreResponse) (string, error) {
	for _, score := range scores {
		p.addScore(score)
	}

	sb := strings.Builder{}

	for i := 0; i < p.Games; i += p.gamesPerLine {
		var limit int = p.gamesPerLine
		if p.Games < i+p.gamesPerLine {
			limit = (i + p.gamesPerLine) - p.Games + 1
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
	return sb.String(), nil
}
