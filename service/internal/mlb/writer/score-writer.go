package writer

import (
	"bytes"
	"embed"
	"github.com/rmarken5/mini-score/service/internal/mlb/fetcher"
	"strings"
	"text/template"
)

//go:embed templates/game-day.gotmpl
var fs embed.FS

type (
	Write interface {
		Write(scoreResponses []*fetcher.FetchScoreResponse) (string, error)
	}

	Writer struct {
	}

	MobileWriter struct {
	}
)

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(scoreResponses []*fetcher.FetchScoreResponse) (string, error) {
	sliceOfBA, err := getSliceOfBA(scoreResponses)
	if err != nil {
		return "", err
	}
	if sliceOfBA == nil {
		return "", nil
	}
	builder := strings.Builder{}

	for i := 1; i+1 < len(sliceOfBA); i++ {
		builder.Write(sliceOfBA[i])
		builder.Write([]byte("\n"))
	}
	return builder.String(), nil
}

func getSliceOfBA(scoreResponses []*fetcher.FetchScoreResponse) ([][]byte, error) {
	if scoreResponses == nil {
		return nil, nil
	}

	sliceOfBA := make([][]byte, len(scoreResponses))

	for i, resp := range scoreResponses {
		b, err := getBytes(resp)
		if err != nil {
			return nil, err
		}
		sliceOfBA[i] = b
	}
	return sliceOfBA, nil
}

func getBytes(scoreResponse *fetcher.FetchScoreResponse) ([]byte, error) {
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
