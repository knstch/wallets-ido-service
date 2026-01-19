package testhelper

import "gorm.io/gorm"

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

	// Postgres: wipe all user-related tables for a clean slate between tests.
	// Note: RESTART IDENTITY makes BIGSERIAL deterministic across tests.
	return s.db.Exec("TRUNCATE TABLE access_tokens, wallets RESTART IDENTITY CASCADE").Error
}
