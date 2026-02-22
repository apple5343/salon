package postgres

import (
	"salon/internal/repository"

	"github.com/jmoiron/sqlx"
)

type modelRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.ModelRepository {
	return &modelRepository{
		db: db,
	}
}
