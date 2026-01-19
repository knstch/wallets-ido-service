package filters

import (
	"gorm.io/gorm"

	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets/models"
)

type WalletsFilter struct {
	ID         uint
	UserID     uint
	Pubkey     string
	Provider   enum.Provider
	IsVerified *bool
}

func BoolPtr(v bool) *bool { return &v }

func (w *WalletsFilter) ToScope() func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Model(&models.UserWallets{})

		if w.ID != 0 {
			tx = tx.Where("id = ?", w.ID)
		}

		if w.UserID != 0 {
			tx = tx.Where("user_id = ?", w.UserID)
		}

		if w.Provider != "" {
			tx = tx.Where("provider = ?", w.Provider.String())
		}

		if w.Pubkey != "" {
			tx = tx.Where("pubkey = ?", w.Pubkey)
		}

		if w.IsVerified != nil {
			if *w.IsVerified {
				tx = tx.Where("verified_at IS NOT NULL")
			} else {
				tx = tx.Where("verified_at IS NULL")
			}
		}

		return tx
	}
}
