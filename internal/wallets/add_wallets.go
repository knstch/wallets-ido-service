package wallets

import (
	"context"

	"github.com/knstch/knstch-libs/tracing"

	"wallets-service/internal/metrics"
	"wallets-service/internal/wallets/filters"
	"wallets-service/internal/wallets/repo"
)

// AddWallets adds multiple wallets for the user.
//
// Behavior:
//   - If a wallet already exists for the same user_id, it is skipped.
//   - If a wallet exists for a different user_id, it is added (same pubkey can belong to multiple users).
//   - All valid wallets are added in a single transaction.
func (s *ServiceImpl) AddWallets(ctx context.Context, userID uint, pubkeys []string) error {
	defer metrics.IncAddWallet()

	ctx, span := tracing.StartSpan(ctx, "wallets: AddWallets")
	defer span.End()

	if len(pubkeys) == 0 {
		return nil
	}

	return s.repo.Transaction(func(st repo.Repository) error {
		for _, pubkey := range pubkeys {
			// Check if wallet already exists for this user
			_, err := st.GetWallet(ctx, filters.WalletsFilter{
				UserID: userID,
				Pubkey: pubkey,
			})
			if err == nil {
				// Wallet already exists for this user, skip it
				continue
			}

			// Wallet doesn't exist for this user, create it
			if err := st.CreateWallet(ctx, userID, pubkey); err != nil {
				// Log error but continue with other wallets
				s.lg.Error("failed to create wallet", err)
				// Continue to next wallet
				continue
			}
		}

		return nil
	})
}
