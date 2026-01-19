package testhelper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"

	_ "wallets-service/migrations"
)

// RunMigrations applies all goose migrations against the provided database.
func RunMigrations(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}

	raw, err := db.DB()
	if err != nil {
		return fmt.Errorf("db.DB: %w", err)
	}

	// Make goose use the same transaction semantics as our app and keep it deterministic.
	goose.SetDialect("postgres")
	goose.SetTableName("goose_db_version")

	migrationsDir, err := findMigrationsDir()
	if err != nil {
		return err
	}

	ctx := context.Background()
	// Ensure connection is valid before running migrations (helps error messages).
	if err := raw.PingContext(ctx); err != nil {
		return fmt.Errorf("db.Ping: %w", err)
	}

	if err := goose.UpContext(ctx, raw, migrationsDir); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}

func findMigrationsDir() (string, error) {
	// Start from this file location and walk upwards until we find go.mod, then use <root>/migrations.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("runtime.Caller failed")
	}

	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, "migrations"), nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find repo root (go.mod) to locate migrations dir")
}
