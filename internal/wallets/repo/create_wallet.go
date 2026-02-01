package repo

import (
	"context"
	"fmt"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/knstch/knstch-libs/tracing"

	"wallets-service/internal/wallets/models"
)

// CreateWallet creates a new wallet row.
//
// If a uniqueness constraint is violated (same user_id + pubkey), CreateWallet returns an error wrapping svcerrs.ErrConflict.
func (r *DBRepo) CreateWallet(ctx context.Context, userID uint, pubkey string) error {
	ctx, span := tracing.StartSpan(ctx, "repo: CreateWallet")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(&models.UserWallets{
		UserID: userID,
		Pubkey: pubkey,
	}).Error; err != nil {
		if isUniqueViolation(err) {
			// This means wallet already exists for this user_id + pubkey combination
			return fmt.Errorf("db.Create: %w", svcerrs.ErrConflict)
		}
		return fmt.Errorf("db.Create: %w", err)
	}

	return nil
}
