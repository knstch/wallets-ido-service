package repo

import (
	"fmt"

	"github.com/knstch/knstch-libs/log"
	"gorm.io/gorm"
)

// NewDBRepo constructs a DB-backed Repository implementation.
func NewDBRepo(lg *log.Logger, db *gorm.DB) (*DBRepo, error) {
	if lg == nil {
		return nil, fmt.Errorf("nil logger")
	}
	if db == nil {
		return nil, fmt.Errorf("nil db")
	}
	return &DBRepo{
		lg: lg,
		db: db,
	}, nil
}
