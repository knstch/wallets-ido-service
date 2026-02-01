package repo

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/tracing"

	"wallets-service/internal/domain/dto"
	"wallets-service/internal/wallets/filters"
	"wallets-service/internal/wallets/models"
)

// GetWallets returns all wallets matching the provided filters.
func (r *DBRepo) GetWallets(ctx context.Context, filters filters.WalletsFilter) ([]dto.Wallet, error) {
	ctx, span := tracing.StartSpan(ctx, "repo: GetWallets")
	defer span.End()

	var wallets []models.UserWallets
	if err := r.db.WithContext(ctx).Scopes(filters.ToScope()).Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("db.Find: %w", err)
	}

	result := make([]dto.Wallet, 0, len(wallets))
	for _, wallet := range wallets {
		result = append(result, dto.Wallet{
			ID:     wallet.ID,
			UserID: wallet.UserID,
			Pubkey: wallet.Pubkey,
		})
	}

	return result, nil
}
