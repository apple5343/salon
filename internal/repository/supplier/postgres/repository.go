package postgres

import (
	"salon/internal/repository"

	"github.com/jmoiron/sqlx"
)

type supplierRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.SupplierRepository {
	return &supplierRepository{
		db: db,
	}
}
