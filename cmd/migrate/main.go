package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"wallets-service/config"

	_ "wallets-service/migrations"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	args := os.Args

	dir, err := filepath.Abs(filepath.Dir(args[0]))
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}

	if err := config.InitENV(dir); err != nil {
		return fmt.Errorf("config.InitENV: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("config.GetConfig: %w", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open: %w", err)
	}

	rawDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("gorm.DB: %w", err)
	}

	if err = goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	flags := flag.NewFlagSet("goose", flag.ExitOnError)
	if len(args) < 2 {
		return fmt.Errorf("no args found")
	}

	if err = flags.Parse(args[1:]); err != nil {
		return fmt.Errorf("flags.Parse: %w", err)
	}

	flagsArgs := flags.Args()

	command := ""
	if len(flagsArgs) >= 2 {
		command = flagsArgs[1]
	} else {
		command = flagsArgs[0]
	}

	ctx := context.Background()
	migrationsDir := filepath.Join(dir, "migrations")
	if err = goose.RunContext(ctx, command, rawDB, migrationsDir, flagsArgs...); err != nil {
		return fmt.Errorf("goose.RunContext: %w", err)
	}

	return nil
}
