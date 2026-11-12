package main

import (
	"embed"
	"log/slog"
	"os"
	"to-do/app"

	_ "to-do/migrations"
)

//go:embed templates/*.html
var templates embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	a := app.New(logger, templates)

	if err := a.Pb.Execute(); err != nil {
		slog.Error("failed to serve", err)
	}
}
