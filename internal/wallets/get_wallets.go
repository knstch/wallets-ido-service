package wallets

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"

	"wallets-service/internal/domain/dto"
	"wallets-service/internal/wallets/filters"
)

// GetWallets returns all wallets associated with the given user ID.
func (s *ServiceImpl) GetWallets(ctx context.Context, userID uint) ([]dto.Wallet, error) {
	ctx, span := tracing.StartSpan(ctx, "wallets: GetWallets")
	defer span.End()

	wallets, err := s.repo.GetWallets(ctx, filters.WalletsFilter{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("repo.GetWallets: %w", err)
	}

	return wallets, nil
}
