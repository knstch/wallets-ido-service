package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"
	"gorm.io/gorm"

	"wallets-service/internal/domain/dto"
	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets/filters"
	"wallets-service/internal/wallets/models"
)

func (r *DBRepo) GetWallet(ctx context.Context, filters filters.WalletsFilter) (dto.Wallet, error) {
	ctx, span := tracing.StartSpan(ctx, "repo: GetWallet")
	defer span.End()

	var wallet models.UserWallets
	if err := r.db.WithContext(ctx).Scopes(filters.ToScope()).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.Wallet{}, fmt.Errorf("wallet not found: %w", svcerrs.ErrDataNotFound)
		}
		return dto.Wallet{}, err
	}

	provider, err := enum.GetProvider(wallet.Provider)
	if err != nil {
		return dto.Wallet{}, fmt.Errorf("enum.GetProvider: %w", err)
	}

	return dto.Wallet{
		ID:         wallet.ID,
		UserID:     wallet.UserID,
		Pubkey:     wallet.Pubkey,
		Provider:   provider,
		VerifiedAt: wallet.VerifiedAt,
	}, nil
}
