package postgres

import (
	"salon/internal/repository"

	"github.com/jmoiron/sqlx"
)

type clientRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.ClientRepository {
	return &clientRepository{
		db: db,
	}
}
