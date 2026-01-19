package testhelper

import "gorm.io/gorm"

// Cleaner is a small helper that resets DB state between integration tests.
//
// It is intentionally best-effort and should only be used in tests.
type Cleaner struct {
	db *gorm.DB
}

// NewCleaner creates a Cleaner instance.
func NewCleaner(db *gorm.DB) Cleaner {
	return Cleaner{db: db}
}

// Clean removes data from test tables.
func (s *Cleaner) Clean() error {
	if s.db == nil {
		return nil
	}

	// Postgres: wipe all service tables for a clean slate between tests.
	// Note: RESTART IDENTITY makes BIGSERIAL deterministic across tests.
	return s.db.Exec("TRUNCATE TABLE user_wallets RESTART IDENTITY CASCADE").Error
}
