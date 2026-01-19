package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/knstch/knstch-libs/log"
	"gorm.io/gorm"

	"wallets-service/internal/domain/dto"
	"wallets-service/internal/domain/enum"
	"wallets-service/internal/wallets/filters"
)

type DBRepo struct {
	lg *log.Logger
	db *gorm.DB
}

type Repository interface {
	// Transaction runs fn inside a database transaction.
	Transaction(fn func(st Repository) error) error
	CreateWallet(ctx context.Context, userID uint, pubkey string, provider enum.Provider) error
	GetWallet(ctx context.Context, filters filters.WalletsFilter) (dto.Wallet, error)
	VerifyWallet(ctx context.Context, filter filters.WalletsFilter) error
}

func (r *DBRepo) NewDBRepo(db *gorm.DB) *DBRepo {
	if db == nil {
		db = r.db.Session(&gorm.Session{NewDB: true})
	}
	return &DBRepo{
		db: db,
		lg: r.lg,
	}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (r *DBRepo) Transaction(fn func(st Repository) error) error {
	db := r.db.Session(&gorm.Session{NewDB: true})
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := fn(r.NewDBRepo(tx)); err != nil {
			return fmt.Errorf("fn: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("db.Transaction: %w", err)
	}
	return nil
}
