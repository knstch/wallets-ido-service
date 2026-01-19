package wallets

import (
	"context"

	"github.com/knstch/knstch-libs/log"
	"github.com/redis/go-redis/v9"

	"wallets-service/config"
	"wallets-service/internal/domain/dto"
	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets/repo"
)

// ServiceImpl is the concrete implementation of Service.
type ServiceImpl struct {
	lg *log.Logger

	repo  repo.Repository
	redis *redis.Client

	cfg config.Config
}

// Service describes the business operations for managing user wallets.
//
// The service persists wallets in Postgres and stores verification challenges in Redis.
type Service interface {
	// AddWallet creates a wallet record for the user (or re-issues a challenge for an existing unverified wallet)
	// and returns a challenge that must be signed to verify ownership.
	AddWallet(ctx context.Context, userID uint, pubkey string, provider enum.Provider) (dto.ChallengeForUser, error)
	// VerifyWallet verifies a previously issued challenge signature and marks the wallet as verified.
	VerifyWallet(ctx context.Context, userID uint, challengeID, signature, pubkey string) error
	// UnlinkWallet removes a wallet record belonging to the user.
	UnlinkWallet(ctx context.Context, walletID, userID uint) error
	// GetWallet returns the wallet for the given user.
	GetWallet(ctx context.Context, userID uint) (dto.Wallet, error)
}

// NewService constructs a wallets service instance.
func NewService(
	lg *log.Logger,
	repo repo.Repository,
	cfg config.Config,
	redis *redis.Client,
) *ServiceImpl {
	return &ServiceImpl{
		lg:    lg,
		repo:  repo,
		cfg:   cfg,
		redis: redis,
	}
}
