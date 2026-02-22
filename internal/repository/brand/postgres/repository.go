package postgres

import (
	"salon/internal/repository"

	"github.com/jmoiron/sqlx"
)

type brandRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.BrandRepository {
	return &brandRepository{
		db: db,
	}
}
