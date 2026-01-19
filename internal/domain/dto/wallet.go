package dto

import (
	"time"

	"wallets-service/internal/domain/enum"
)

type Wallet struct {
	ID         uint
	UserID     uint
	Pubkey     string
	Provider   enum.Provider
	VerifiedAt *time.Time
}
