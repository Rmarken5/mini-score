package handlers

import (
	"embed"
	"github.com/labstack/echo/v4"
	"html/template"
	"log"
)

//go:embed templates/index.gotmpl
var files embed.FS

type (
	IndexHandler struct {
		logger *log.Logger
	}
)

func NewIndexHandler(logger *log.Logger) *IndexHandler {
	return &IndexHandler{logger: logger}
}

func (h *IndexHandler) ServeHTTP(c echo.Context) error {

	file, err := template.ParseFS(files, "templates/index.gotmpl")
	if err != nil {
		h.logger.Printf("Error opening file from template: %s", err)
		return err
	}

	return file.Execute(c.Response(), nil)
}
