package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"
	"gorm.io/gorm"

	"wallets-service/internal/wallets/filters"
	"wallets-service/internal/wallets/models"
)

func (r *DBRepo) VerifyWallet(ctx context.Context, filter filters.WalletsFilter) error {
	ctx, span := tracing.StartSpan(ctx, "repo: VerifyWallet")
	defer span.End()

	var walletFromDB models.UserWallets
	if err := r.db.WithContext(ctx).Scopes(filter.ToScope()).First(&walletFromDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("wallet to update is not found: %w", svcerrs.ErrDataNotFound)
		}
		return fmt.Errorf("db.First: %w", err)
	}

	now := time.Now()
	walletFromDB.VerifiedAt = &now

	if err := r.db.Save(&walletFromDB).Error; err != nil {
		return fmt.Errorf("db.Save: %w", err)
	}

	return nil
}
