package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitUserWalletsTable, downInitUserWalletsTable)
}

func upInitUserWalletsTable(_ context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`
			CREATE TABLE user_wallets (
			  id BIGSERIAL PRIMARY KEY,
			  user_id BIGINT NOT NULL,
			  pubkey TEXT NOT NULL,
			  provider TEXT NOT NULL,
			  verified_at TIMESTAMPTZ,
			  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			  UNIQUE(pubkey),
			  UNIQUE(user_id, pubkey)
			);
`); err != nil {
		return err
	}
	return nil
}

func downInitUserWalletsTable(_ context.Context, _ *sql.Tx) error {
	return nil
}
