package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"
	"gorm.io/gorm"

	"wallets-service/internal/wallets/filters"
	"wallets-service/internal/wallets/models"
)

// DeleteWallet deletes a wallet matching the provided filters.
//
// If no wallet matches, DeleteWallet returns an error wrapping svcerrs.ErrDataNotFound.
func (r *DBRepo) DeleteWallet(ctx context.Context, filters filters.WalletsFilter) error {
	ctx, span := tracing.StartSpan(ctx, "repo: DeleteWallet")
	defer span.End()

	var wallet models.UserWallets
	if err := r.db.WithContext(ctx).Scopes(filters.ToScope()).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("wallet not found: %w", svcerrs.ErrDataNotFound)
		}
		return fmt.Errorf("db.First: %w", err)
	}

	if err := r.db.WithContext(ctx).Delete(&wallet, wallet.ID).Error; err != nil {
		return fmt.Errorf("db.Delete: %w", err)
	}

	return nil
}
