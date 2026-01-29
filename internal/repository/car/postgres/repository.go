package postgres

import (
	"salon/internal/repository"

	"github.com/jmoiron/sqlx"
)

type carRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.CarRepository {
	return &carRepository{
		db: db,
	}
}
