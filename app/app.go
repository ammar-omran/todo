package app

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type App struct {
	logger    *slog.Logger
	pb        *pocketbase.PocketBase
	templates fs.FS
}

func New(logger *slog.Logger, templates fs.FS) *App {
	pb := pocketbase.NewWithConfig(pocketbase.Config{
		HideStartBanner: true,
		DefaultDataDir:  "./pb_data",
	})

	return &App{
		logger:    logger,
		pb:        pb,
		templates: templates,
	}
}

func (a *App) Start(ctx context.Context) error {
	tmpl := template.Must(template.New("").ParseFS(a.templates, "templates/*"))

	a.loadRoutes(tmpl)

	if err := a.pb.Bootstrap(); err != nil {
		return fmt.Errorf("bootstrap pocketbase: %w", err)
	}
	defer a.pb.ClearBootstrap()

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			a.logger.Info("shutting down")
			event := &core.TerminateEvent{App: a.pb}
			err := a.pb.OnTerminate().Trigger(event, func(e *core.TerminateEvent) error {
				return e.Next()
			})
			if err != nil {
				a.logger.Error("failed to terminate app", slog.Any("error", err))
			}
		case <-done:
		}
	}()

	a.logger.Info("server listening", slog.String("addr", ":8080"))

	err := apis.Serve(a.pb, apis.ServeConfig{
		HttpAddr:        ":8080",
		ShowStartBanner: false,
	})
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}
