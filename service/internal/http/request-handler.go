package http

import (
	"github.com/labstack/echo/v4"
	"github.com/rmarken5/mini-score/service/internal/mlb/facade"
	"net/http"
	"time"
)

type (
	Server struct {
		facade facade.ScoreFacade
	}
)

const layout = "2006-01-02"

func NewServer(facade facade.ScoreFacade) *Server {
	return &Server{facade: facade}
}

func (s *Server) PrintGames(c echo.Context) error {
	date := c.Param("date")
	if date == "" {
		date = time.Now().Format(layout)
	}

	dateObj, err := time.Parse(layout, date)
	if err != nil {
		return err
	}

	scores, err := s.facade.ProcessScores(dateObj)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, scores)

}
