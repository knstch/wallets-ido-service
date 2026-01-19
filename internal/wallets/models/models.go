package models

import "time"

type UserWallets struct {
	ID         uint
	UserID     uint
	Pubkey     string
	Provider   string
	VerifiedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (UserWallets) TableName() string {
	return "user_wallets"
}
