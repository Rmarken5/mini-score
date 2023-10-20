package rest

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/data-access/db/repository"
	"github.com/rmarken5/mini-score/service/internal/nfl/logic/internal/general"
	"github.com/rs/zerolog"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	gameTimeFormat    = "Mon, 3:04 PM"
	headingTimeFormat = "Jan, 02 2006"
)

var _ ScoreboardFacade = &Controller{}

type (
	ScoreboardFacade interface {
		GetScoreboardForDate(date time.Time) (Scores, error)
	}

	Controller struct {
		logger zerolog.Logger
		repo   repository.Repository
	}

	score struct {
		gameID             string
		awayTeam, homeTeam team
		quarter            string
		gameClock          string
		startTime          time.Time
	}
	team struct {
		name   string
		scores []string
	}
	Scores []score

	ByGameTime Scores
)

func NewScoreboardFacade(logger zerolog.Logger, db *sqlx.DB) *Controller {
	return &Controller{
		logger: logger,
		repo:   repository.NewRepository(logger, db),
	}
}

func (c *Controller) GetScoreboardForDate(date time.Time) (Scores, error) {
	logger := c.logger.With().Str("method", "GetScoreboardForDate").Logger()
	logger.Info().Msgf("getting scores at %s", date)

	start := general.StartTime(date)
	end := general.EndTime(date)

	gqs, err := c.repo.GetGameTeamQuarterScore(start, &end)
	if err != nil {
		logger.Error().Err(err).Msg("while getting scores")
	}

	games, err := c.repo.GetGamesWithTeamAbv(start, &end)
	if err != nil {
		logger.Error().Err(err).Msg("while getting games")
	}

	scores := buildScoreFromDB(gqs, games)
	sort.Sort(ByGameTime(scores))

	return scores, nil
}

func buildScoreFromDB(gts []repository.GameTeamQuarterScore, games []repository.Game) Scores {

	scores := make(Scores, 0)
	for _, g := range games {
		gameClock := g.GameClock
		if gameClock == "" {
			gameClock = g.GameTime.Local().Format(gameTimeFormat)
		}
		s := score{
			gameID: g.ID,
			awayTeam: team{
				name:   g.AwayTeam,
				scores: getScoresForTeam(g.ID, g.AwayTeam, gts),
			},
			homeTeam: team{
				name:   g.HomeTeam,
				scores: getScoresForTeam(g.ID, g.HomeTeam, gts),
			},
			quarter:   g.Quarter,
			gameClock: gameClock,
			startTime: g.GameTime,
		}

		scores = append(scores, s)
	}
	return scores

}

func getScoresForTeam(gameID, teamName string, quarterScores []repository.GameTeamQuarterScore) []string {
	scores := make([]string, 0)

	for _, qs := range quarterScores {
		if qs.GameID == gameID && qs.TeamAbbreviation == teamName {
			scores = append(scores, qs.Score)
		}
	}
	return scores
}

func (s Scores) PrintScoreboard(writer io.Writer, scoresPerLine int) error {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s\n", time.Now().Local().Format(headingTimeFormat)))
	for i := 0; i < len(s); {
		diff := 0
		if scoresPerLine+i >= len(s) {
			diff = (scoresPerLine + i) - len(s)
		}
		upperBounds := (scoresPerLine + i) - diff

		if i > 0 {
			sb.WriteString("\n")
		}
		for _, score := range s[i:upperBounds] {
			sb.WriteString(score.buildTopAndBottomBoarder() + " ")
		}
		sb.WriteString("\n")
		for _, score := range s[i:upperBounds] {
			sb.WriteString(score.buildQuarterLine() + " ")
		}
		sb.WriteString("\n")
		for _, score := range s[i:upperBounds] {
			sb.WriteString(score.buildAwayLine() + " ")
		}
		sb.WriteString("\n")
		for _, score := range s[i:upperBounds] {
			sb.WriteString(score.buildGameClockLine() + " ")
		}
		sb.WriteString("\n")
		for _, score := range s[i:upperBounds] {
			sb.WriteString(score.buildHomeLine() + " ")
		}
		sb.WriteString("\n")
		for _, score := range s[i:upperBounds] {
			sb.WriteString(score.buildTopAndBottomBoarder() + " ")
		}
		i += scoresPerLine
	}

	_, err := writer.Write([]byte(sb.String()))
	if err != nil {
		return err
	}
	return nil
}

func (s score) buildQuarterLine() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("* %-3s", "Q"))
	for i, _ := range s.awayTeam.scores {
		sb.WriteString(fmt.Sprintf("%3d", i+1))
	}
	sb.WriteString(fmt.Sprintf("%5s", ""))
	var gapSize = 2
	if len(s.awayTeam.scores)%2 == 0 {
		gapSize = 3
	}
	sb.WriteString(fmt.Sprintf("%"+strconv.Itoa(gapSize)+"s", "*"))
	return sb.String()
}

func (s score) buildAwayLine() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("* %-3s", s.awayTeam.name))
	var currentScore int
	for _, val := range s.awayTeam.scores {
		iVal, err := strconv.Atoi(val)
		if err != nil {
			iVal = 0
			// TODO log
		}
		currentScore += iVal
		sb.WriteString(fmt.Sprintf("%3d", iVal))
	}
	sb.WriteString(fmt.Sprintf("%5d", currentScore))
	var gapSize = 2
	if len(s.awayTeam.scores)%2 == 0 {
		gapSize = 3
	}
	sb.WriteString(fmt.Sprintf("%"+strconv.Itoa(gapSize)+"s", "*"))
	return sb.String()
}

func (s score) buildHomeLine() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("* %-3s", s.homeTeam.name))
	var currentScore int
	for _, val := range s.homeTeam.scores {
		iVal, err := strconv.Atoi(val)
		if err != nil {
			iVal = 0
			// TODO log
		}
		currentScore += iVal
		sb.WriteString(fmt.Sprintf("%3d", iVal))
	}
	sb.WriteString(fmt.Sprintf("%5d", currentScore))
	var gapSize = 2
	if len(s.homeTeam.scores)%2 == 0 {
		gapSize = 3
	}
	sb.WriteString(fmt.Sprintf("%"+strconv.Itoa(gapSize)+"s", "*"))
	return sb.String()
}

func (s score) buildGameClockLine() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("* Q%-2s", s.quarter))

	lineLen := 3*len(s.homeTeam.scores) + 5
	sb.WriteString(fmt.Sprintf("%"+strconv.Itoa(lineLen)+"s", s.gameClock))

	var gapSize = 2
	if len(s.homeTeam.scores)%2 == 0 {
		gapSize = 3
	}
	sb.WriteString(fmt.Sprintf("%"+strconv.Itoa(gapSize)+"s", "*"))

	return sb.String()
}

func (s score) buildTopAndBottomBoarder() string {
	sb := strings.Builder{}
	lineLen := (10 + 3*len(s.awayTeam.scores)) / 2
	for i := 0; i < lineLen; i++ {
		sb.WriteString("* ")
	}

	sb.WriteString("* *")

	return sb.String()
}

func (b ByGameTime) Len() int {
	return len(b)
}

func (b ByGameTime) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (b ByGameTime) Less(i, j int) bool {
	if b[i].startTime.Before(b[j].startTime) {
		return true
	} else if b[j].startTime.Before(b[i].startTime) {
		return false
	}
	return b[i].gameID < b[j].gameID
}
