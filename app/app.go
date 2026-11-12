package app

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"

	"ganigeorgiev.com/pocketbase/pocketbase"
	"ganigeorgiev.com/pocketbase/pocketbase/apis"
)

type App struct {
	logger    *slog.Logger
	Pb        *pocketbase.PocketBase
	templates fs.FS
}

func New(logger *slog.Logger, templates fs.FS) *App {
	pb := pocketbase.NewWithConfig(pocketbase.Config{
		HideStartBanner: false,
		DefaultDataDir:  "./pb_data",
	})

	return &App{
		logger:    logger,
		Pb:        pb,
		templates: templates,
	}
}

func (a *App) Start(ctx context.Context) error {
	tmpl := template.Must(template.New("").ParseFS(a.templates, "templates/*"))

	a.loadRoutes(tmpl)

	if err := a.Pb.Bootstrap(); err != nil {
		return fmt.Errorf("bootstrap pocketbase: %w", err)
	}
	defer a.Pb.ClearBootstrap()

	a.logger.Info("server listening", slog.String("addr", "localhost:8080"))

	err := apis.Serve(a.Pb, apis.ServeConfig{
		HttpAddr:        "localhost:8080",
		ShowStartBanner: true,
	})
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}
