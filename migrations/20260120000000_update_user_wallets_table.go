package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdateUserWalletsTable, downUpdateUserWalletsTable)
}

func upUpdateUserWalletsTable(_ context.Context, tx *sql.Tx) error {
	// Remove verified_at column if it exists
	if _, err := tx.Exec(`
		ALTER TABLE user_wallets 
		DROP COLUMN IF EXISTS verified_at;
	`); err != nil {
		return err
	}

	// Remove provider column if it exists
	if _, err := tx.Exec(`
		ALTER TABLE user_wallets 
		DROP COLUMN IF EXISTS provider;
	`); err != nil {
		return err
	}

	// Drop old unique constraint on pubkey if exists
	if _, err := tx.Exec(`
		DROP INDEX IF EXISTS user_wallets_pubkey_key;
	`); err != nil {
		return err
	}

	// Add unique constraint on (user_id, pubkey)
	if _, err := tx.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS user_wallets_user_id_pubkey_unique 
		ON user_wallets(user_id, pubkey) 
		WHERE deleted_at IS NULL;
	`); err != nil {
		return err
	}

	return nil
}

func downUpdateUserWalletsTable(_ context.Context, tx *sql.Tx) error {
	// Revert changes
	if _, err := tx.Exec(`
		DROP INDEX IF EXISTS user_wallets_user_id_pubkey_unique;
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS user_wallets_pubkey_key 
		ON user_wallets(pubkey) 
		WHERE deleted_at IS NULL;
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		ALTER TABLE user_wallets 
		ADD COLUMN IF NOT EXISTS provider TEXT NOT NULL DEFAULT 'PHANTOM';
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		ALTER TABLE user_wallets 
		ADD COLUMN IF NOT EXISTS verified_at TIMESTAMPTZ NULL;
	`); err != nil {
		return err
	}

	return nil
}
