package main

import (
	"context"
	"errors"
	"fmt"
	defaultLog "log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/knstch/knstch-libs/endpoints"
	"github.com/knstch/knstch-libs/log"
	"github.com/knstch/knstch-libs/tracing"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"wallets-service/config"
	"wallets-service/internal/endpoints/public"
	"wallets-service/internal/wallets"
	"wallets-service/internal/wallets/repo"
)

func main() {
	if err := run(); err != nil {
		defaultLog.Println(err)
		recover()
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

	shutdown := tracing.InitTracer(cfg.ServiceName, cfg.JaegerHost)
	defer shutdown(context.Background())

	logger := log.NewLogger(cfg.ServiceName, log.InfoLevel)

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open: %w", err)
	}

	dbRepo, err := repo.NewDBRepo(logger, db)
	if err != nil {
		return fmt.Errorf("repo.NewDBRepo: %w", err)
	}

	dsnRedis, err := redis.ParseURL(cfg.GetRedisDSN())
	if err != nil {
		return err
	}
	redisClient := redis.NewClient(dsnRedis)

	svc := wallets.NewService(logger, dbRepo, *cfg, redisClient)

	publicController := public.NewController(svc, logger, cfg)
	publicEndpoints := endpoints.InitHttpEndpoints(cfg.ServiceName, publicController.Endpoints())

	srv := http.Server{
		Addr: ":" + cfg.PublicHTTPAddr,
		Handler: http.TimeoutHandler(
			publicEndpoints,
			time.Second*5,
			"service temporary unavailable",
		),
		ReadHeaderTimeout: time.Millisecond * 500,
		ReadTimeout:       time.Minute * 5,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err = srv.Shutdown(context.Background()); err != nil {
			logger.Error("error shutting down", err)
		}
		close(idleConnsClosed)
	}()

	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Error("error serving", err)
	}

	<-idleConnsClosed

	return nil
}
