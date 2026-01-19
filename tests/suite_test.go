package wallets_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	knlog "github.com/knstch/knstch-libs/log"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"wallets-service/config"
	"wallets-service/internal/wallets"
	"wallets-service/internal/wallets/repo"
	"wallets-service/testhelper"
)

func TestWalletsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(WalletsServiceTestSuite))
}

type WalletsServiceTestSuite struct {
	suite.Suite

	cfg config.Config
	svc wallets.Service

	db      *gorm.DB
	rdb     *redis.Client
	cleaner testhelper.Cleaner

	logger *knlog.Logger
	dbRepo repo.Repository
}

func (s *WalletsServiceTestSuite) SetupSuite() {
	t := s.Require()
	time.Local = time.UTC

	// Tests should always prefer test creds; only load test.env.
	root := mustFindRepoRoot(t)
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
	if cfg.RedisConfig.Host == "" {
		cfg.RedisConfig.Host = envOr("REDIS_HOST", "localhost")
	}
	if cfg.RedisConfig.Port == "" {
		cfg.RedisConfig.Port = envOr("REDIS_PORT", "6379")
	}
	s.cfg = cfg

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	t.NoError(err)
	s.db = db

	t.NoError(testhelper.RunMigrations(db))
	s.cleaner = testhelper.NewCleaner(db)

	dsnRedis, err := redis.ParseURL(cfg.GetRedisDSN())
	t.NoError(err)
	s.rdb = redis.NewClient(dsnRedis)
	t.NoError(s.rdb.Ping(context.Background()).Err())

	logger := newTestLogger(cfg.ServiceName)
	dbRepo, err := repo.NewDBRepo(logger, db)
	t.NoError(err)

	s.logger = logger
	s.dbRepo = dbRepo
	s.svc = wallets.NewService(logger, dbRepo, cfg, s.rdb)
}

func (s *WalletsServiceTestSuite) SetupTest() {
	s.Require().NoError(s.cleaner.Clean())
	s.Require().NoError(s.rdb.FlushDB(context.Background()).Err())
}
