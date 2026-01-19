package wallets

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"

	"wallets-service/internal/wallets/filters"
)

// UnlinkWallet deletes a wallet by ID for the given user.
//
// If the wallet does not exist or does not belong to the user, UnlinkWallet returns an error wrapping svcerrs.ErrDataNotFound.
func (s *ServiceImpl) UnlinkWallet(ctx context.Context, walletID, userID uint) error {
	ctx, span := tracing.StartSpan(ctx, "wallets: UnlinkWallet")
	defer span.End()

	if err := s.repo.DeleteWallet(ctx, filters.WalletsFilter{
		ID:     walletID,
		UserID: userID,
	}); err != nil {
		return fmt.Errorf("repo.DeleteWallet: %w", err)
	}

	return nil
}
