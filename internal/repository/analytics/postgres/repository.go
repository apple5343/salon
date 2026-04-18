package postgres

import (
	"salon/internal/repository"

	"github.com/jmoiron/sqlx"
)

type analyticsRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.AnalyticsRepository {
	return &analyticsRepository{
		db: db,
	}
}
