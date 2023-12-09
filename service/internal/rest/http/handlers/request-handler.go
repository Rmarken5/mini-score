package handlers

import (
	"github.com/labstack/echo/v4"
	mlbfacade "github.com/rmarken5/mini-score/service/internal/mlb/facade"
	nflfacade "github.com/rmarken5/mini-score/service/internal/nfl/logic/rest"
	user_agent "github.com/rmarken5/mini-score/service/internal/rest/user-agent"
	"net/http"
	"time"
)

type (
	Server struct {
		mlbFacade mlbfacade.ScoreFacade
		nflFacade nflfacade.ScoreboardFacade
	}
)

const layout = "2006-01-02"

func NewServer(mlbFacade mlbfacade.ScoreFacade, scoreboard nflfacade.ScoreboardFacade) *Server {
	return &Server{mlbFacade: mlbFacade, nflFacade: scoreboard}
}

func (s *Server) PrintBaseballGames(c echo.Context) error {
	date := c.Param("date")
	if date == "" {
		date = time.Now().Format(layout)
	}

	dateObj, err := time.Parse(layout, date)
	if err != nil {
		return err
	}

	scores, err := mlbfacade.ProcessScores(s.mlbFacade, c.Request().Context(), dateObj)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/plain")
	return c.String(http.StatusOK, scores)

}

func (s *Server) PrintFootballGames(c echo.Context) error {
	date := c.Param("date")
	if date == "" {
		date = time.Now().Format(layout)
	}

	gamesPerLine := 1
	if !user_agent.IsMobile(c.Request().Context()) {
		gamesPerLine = 3
	}

	dateObj, err := time.Parse(layout, date)
	if err != nil {
		return err
	}

	scores, err := s.nflFacade.GetScoreboardForDate(dateObj)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Content-Type", "text/plain")
	err = scores.PrintScoreboard(c.Response(), dateObj, gamesPerLine)
	if err != nil {
		return err
	}
	return nil
}
