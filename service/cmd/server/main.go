package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/rmarken5/mini-score/service/cmd/internal"
	mlbfacade "github.com/rmarken5/mini-score/service/internal/mlb/facade"
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	nflfacade "github.com/rmarken5/mini-score/service/internal/nfl/logic/rest"
	"github.com/rmarken5/mini-score/service/internal/rest/http/handlers"
	agent "github.com/rmarken5/mini-score/service/internal/rest/user-agent"
	"github.com/rs/zerolog"
	"log"
	h "net/http"
	"os"
	"time"
	_ "time/tzdata"
)

func init() {
	timezone := "America/New_York"
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	// Set the default timezone
	time.Local = loc
}

func main() {
	logger := createLogger()
	httpClient := &h.Client{}
	fetch := fetcher.NewFetcher(httpClient)
	mlbFacade := mlbfacade.NewScoreFacadeImpl(fetch, fetch)
	nflFacade := nflfacade.NewScoreboardFacade(logger, internal.MustConnectDatabase(logger))

	s := handlers.NewServer(mlbFacade, nflFacade)

	idxHandler := handlers.NewIndexHandler(&log.Logger{})
	e := echo.New()
	e.Use(agent.HandleUserAgent)
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.GET("/", idxHandler.ServeHTTP)
	e.GET("/mlb/:date", s.PrintBaseballGames)
	e.GET("/mlb", s.PrintBaseballGames)
	e.GET("/nfl/:date", s.PrintFootballGames)
	e.GET("/nfl", s.PrintFootballGames)

	httpServer := h.Server{Addr: ":8080", Handler: e}

	if err := httpServer.ListenAndServe(); !errors.Is(err, h.ErrServerClosed) {
		log.Fatal(err)
	}

}

func createLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.With().Str("service", "rest").Logger()
	return logger
}
