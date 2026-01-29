package postgres

import (
	"salon/internal/repository"

	"github.com/jmoiron/sqlx"
)

type employeeRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}
