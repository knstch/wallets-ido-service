package wallets

import (
	"context"

	"github.com/knstch/knstch-libs/log"

	"wallets-service/config"
	"wallets-service/internal/domain/dto"
	"wallets-service/internal/wallets/repo"
)

// ServiceImpl is the concrete implementation of Service.
type ServiceImpl struct {
	lg *log.Logger

	repo repo.Repository

	cfg config.Config
}

// Service describes the business operations for managing user wallets.
type Service interface {
	// AddWallets adds multiple wallets for the user.
	// If a wallet already exists for the same user_id, it is skipped.
	// If a wallet exists for a different user_id, it is added.
	AddWallets(ctx context.Context, userID uint, pubkeys []string) error
	// GetWallet returns the wallet for the given user.
	GetWallet(ctx context.Context, userID uint) (dto.Wallet, error)
	// GetWallets returns all wallets for the given user.
	GetWallets(ctx context.Context, userID uint) ([]dto.Wallet, error)
}

// NewService constructs a wallets service instance.
func NewService(
	lg *log.Logger,
	repo repo.Repository,
	cfg config.Config,
) *ServiceImpl {
	return &ServiceImpl{
		lg:   lg,
		repo: repo,
		cfg:  cfg,
	}
}
