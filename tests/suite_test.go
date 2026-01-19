package users_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	googlemocks "wallets-service/internal/connector/google/mocks"

	"wallets-service/config"
	"wallets-service/internal/wallets"
	"wallets-service/internal/wallets/repo"
	"wallets-service/testhelper"
)

func TestUsersServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UsersServiceTestSuite))
}

type UsersServiceTestSuite struct {
	suite.Suite

	cfg config.Config
	svc wallets.Service

	db      *gorm.DB
	rdb     *redis.Client
	cleaner testhelper.Cleaner

	googleMock *googlemocks.Client
}

func (s *UsersServiceTestSuite) SetupSuite() {
	t := require.New(s.T())
	time.Local = time.UTC

	// Load .env and then override with test.env (tests should always prefer test creds).
	root := mustFindRepoRoot(t)
	t.NoError(config.InitENV(root))
	err := godotenv.Overload(filepath.Join(root, "test.env"))
	t.NoError(err)

	cfgPtr, err := config.GetConfig()
	t.NoError(err)
	cfg := *cfgPtr

	// Provide safe defaults for local tests if some envs are missing.
	if cfg.ServiceName == "" {
		cfg.ServiceName = "wallets-service-test"
	}
	if cfg.JwtSecret == "" {
		cfg.JwtSecret = "test-secret"
	}
	if cfg.PlatformURL == "" {
		cfg.PlatformURL = "https://example.com"
	}
	// Redis envs in your .env may be named REDIS_EXTERNAL_PORT, so keep a fallback.
	if cfg.RedisConfig.Host == "" {
		cfg.RedisConfig.Host = envOr("REDIS_HOST", "localhost")
	}
	if cfg.RedisConfig.Port == "" {
		cfg.RedisConfig.Port = envOr("REDIS_PORT", envOr("REDIS_EXTERNAL_PORT", "6379"))
	}
	s.cfg = cfg

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		s.T().Skipf("postgres is not available for tests: %v", err)
	}
	s.db = db

	t.NoError(testhelper.RunMigrations(db))
	s.cleaner = testhelper.NewCleaner(db)

	s.rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Host + ":" + cfg.RedisConfig.Port,
		Username: cfg.RedisConfig.Username,
		Password: cfg.RedisConfig.Password,
	})
	if err := s.rdb.Ping(context.Background()).Err(); err != nil {
		s.T().Skipf(
			"redis is not available for tests: %v (addr=%s password_len=%d)",
			err,
			cfg.RedisConfig.Host+":"+cfg.RedisConfig.Port,
			len(cfg.RedisConfig.Password),
		)
	}

	logger := newTestLogger(cfg.ServiceName)
	dbRepo, err := repo.NewDBRepo(logger, db)
	t.NoError(err)

	s.googleMock = &googlemocks.Client{}
	s.svc = wallets.NewService(logger, dbRepo, cfg, s.googleMock, s.rdb)
}

func (s *UsersServiceTestSuite) SetupTest() {
	if err := s.cleaner.Clean(); err != nil {
		panic(err)
	}
	_ = s.rdb.FlushDB(context.Background()).Err()
	s.googleMock.ExpectedCalls = nil
	s.googleMock.Calls = nil
}
