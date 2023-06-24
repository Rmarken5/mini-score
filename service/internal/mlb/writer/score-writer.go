package writer

import (
	"bytes"
	"embed"
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"text/template"
)

//go:embed templates/game-day.gotmpl
var fs embed.FS

type (
	Write interface {
		Write(scoreResponse *fetcher.FetchScoreResponse) ([]byte, error)
	}

	Writer struct {
	}
)

func (w *Writer) Write(scoreResponse *fetcher.FetchScoreResponse) ([]byte, error) {
	file, err := template.ParseFS(fs, "templates/game-day.gotmpl")
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(nil)

	err = file.Execute(buff, scoreResponse)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
