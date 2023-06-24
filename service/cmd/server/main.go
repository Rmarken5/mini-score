package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rmarken5/mini-score/service/internal/http"
	"github.com/rmarken5/mini-score/service/internal/mlb/facade"
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"log"
	h "net/http"
)

func main() {
	httpClient := &h.Client{}
	fetch := fetcher.NewFetcher(httpClient)
	f := facade.NewScoreFacadeImpl(fetch, fetch)
	s := http.NewServer(f)
	e := echo.New()
	e.GET("/mlb/:date", s.PrintGames)

	httpServer := h.Server{Addr: ":8080", Handler: e}

	if err := httpServer.ListenAndServe(); err != h.ErrServerClosed {
		log.Fatal(err)
	}

}
