package filters

import (
	"gorm.io/gorm"

	"wallets-service/internal/wallets/models"
)

// WalletsFilter defines query parameters for selecting wallets.
type WalletsFilter struct {
	ID     uint
	UserID uint
	Pubkey string
}

// ToScope converts the filter to a GORM scope.
func (w *WalletsFilter) ToScope() func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Model(&models.UserWallets{})

		if w.ID != 0 {
			tx = tx.Where("id = ?", w.ID)
		}

		if w.UserID != 0 {
			tx = tx.Where("user_id = ?", w.UserID)
		}

		if w.Pubkey != "" {
			tx = tx.Where("pubkey = ?", w.Pubkey)
		}

		return tx
	}
}
