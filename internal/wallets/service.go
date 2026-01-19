package wallets

import (
	"github.com/knstch/knstch-libs/log"
	"github.com/redis/go-redis/v9"

	"wallets-service/config"
	"wallets-service/internal/wallets/repo"
)

type ServiceImpl struct {
	lg *log.Logger

	repo  repo.Repository
	redis redis.Client

	cfg config.Config
}

type Service interface{}

// NewService constructs the Users service.
func NewService(
	lg *log.Logger,
	repo repo.Repository,
	cfg config.Config,
	redis redis.Client,
) *ServiceImpl {
	return &ServiceImpl{
		lg:    lg,
		repo:  repo,
		cfg:   cfg,
		redis: redis,
	}
}
