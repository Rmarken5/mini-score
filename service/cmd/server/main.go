package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rmarken5/mini-score/service/internal/http"
	"github.com/rmarken5/mini-score/service/internal/mlb/facade"
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"github.com/rmarken5/mini-score/service/internal/mlb/writer"
	"log"
	h "net/http"
)

func main() {
	httpClient := &h.Client{}
	fetch := fetcher.NewFetcher(httpClient)
	write := writer.NewWriter()
	f := facade.NewScoreFacadeImpl(fetch, fetch, write)
	s := http.NewServer(f)
	e := echo.New()
	e.GET("/mlb/:date", s.PrintGames)

	httpServer := h.Server{Addr: ":8080", Handler: e}

	if err := httpServer.ListenAndServe(); err != h.ErrServerClosed {
		log.Fatal(err)
	}

}
