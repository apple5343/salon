package postgres

import (
	"salon/internal/repository"
	"salon/pkg/clock"

	"github.com/jmoiron/sqlx"
)

type saleRepository struct {
	db    *sqlx.DB
	clock clock.Clock
}

func NewRepository(db *sqlx.DB, clock clock.Clock) repository.SaleRepository {
	return &saleRepository{
		db:    db,
		clock: clock,
	}
}
