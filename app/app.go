package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"time"
	"to-do/database"
)

type App struct {
	logger     *slog.Logger
	router     *http.ServeMux
	dbConn     *sql.DB
	migrations fs.FS
	templates  fs.FS
}

func New(logger *slog.Logger, migrations fs.FS, templates fs.FS) *App {
	router := http.NewServeMux()

	app := &App{
		logger:     logger,
		router:     router,
		migrations: migrations,
		templates:  templates,
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	conn, err := database.Connect(ctx, a.logger, a.migrations)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	a.dbConn = conn

	tmpl := template.Must(template.New("").ParseFS(a.templates, "templates/*"))

	a.loadRoutes(tmpl)

	server := http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	// graceful shutdown
	done := make(chan struct{})
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("failed to listen and serve", slog.Any("error", err))
		}
		close(done)

	}()

	a.logger.Info("server listening", slog.String("addr", ":8080"))
	select {
	case <-done:
		break
	case <-ctx.Done():
		a.logger.Info("Shutting Down!")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		server.Shutdown(ctx)
		cancel()
	}

	return nil
}
