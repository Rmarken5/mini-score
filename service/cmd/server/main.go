package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	genhandle "github.com/rmarken5/mini-score/service/internal/general/http/handlers"
	genmid "github.com/rmarken5/mini-score/service/internal/general/user-agent"
	"github.com/rmarken5/mini-score/service/internal/mlb/facade"
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"github.com/rmarken5/mini-score/service/internal/mlb/http/http"
	"log"
	h "net/http"
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

	httpClient := &h.Client{}
	fetch := fetcher.NewFetcher(httpClient)
	f := facade.NewScoreFacadeImpl(fetch, fetch)
	s := http.NewServer(f)
	idxHandler := genhandle.NewIndexHandler(&log.Logger{})
	e := echo.New()
	e.Use(genmid.HandleUserAgent)
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.GET("/", idxHandler.ServeHTTP)
	e.GET("/mlb/:date", s.PrintGames)
	e.GET("/mlb", s.PrintGames)

	httpServer := h.Server{Addr: ":8080", Handler: e}

	if err := httpServer.ListenAndServe(); err != h.ErrServerClosed {
		log.Fatal(err)
	}

}
