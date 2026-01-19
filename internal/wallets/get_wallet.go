package wallets

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"

	"wallets-service/internal/domain/dto"
	"wallets-service/internal/wallets/filters"
)

// GetWallet returns the wallet associated with the given user ID.
//
// If the wallet doesn't exist, GetWallet returns an error wrapping svcerrs.ErrDataNotFound.
func (s *ServiceImpl) GetWallet(ctx context.Context, userID uint) (dto.Wallet, error) {
	ctx, span := tracing.StartSpan(ctx, "wallets: GetWallet")
	defer span.End()

	wallet, err := s.repo.GetWallet(ctx, filters.WalletsFilter{UserID: userID})
	if err != nil {
		return dto.Wallet{}, fmt.Errorf("repo.GetWallet: %w", err)
	}

	return wallet, nil
}
