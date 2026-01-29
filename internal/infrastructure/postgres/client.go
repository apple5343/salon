package postgres

import (
	"salon/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewClient(cfg *config.Postgres) (*sqlx.DB, error) {
	client, err := sqlx.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		return nil, err
	}
	return client, nil
}
