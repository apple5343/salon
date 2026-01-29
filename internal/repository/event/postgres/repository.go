package postgres

import (
	"salon/internal/repository"

	"github.com/jmoiron/sqlx"
)

type eventRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.EventRepository {
	return &eventRepository{
		db: db,
	}
}
