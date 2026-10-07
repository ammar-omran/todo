package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

func Connect(
	ctx context.Context,
	logger *slog.Logger,
	migrations fs.FS,
) (*sql.DB, error) {
	const dbPath = "todo.db"

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Debug("running migrations")

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("create migration source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance(
		"iofs",
		source,
		"sqlite3://todo.db",
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("create migrator: %w", err)
	}
	defer func() {
		srcErr, dbErr := migrator.Close()
		if srcErr != nil {
			logger.Error("failed to close migration source", "error", srcErr)
		}
		if dbErr != nil {
			logger.Error("failed to close migration database", "error", dbErr)
		}
	}()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		conn.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return conn, nil
}
