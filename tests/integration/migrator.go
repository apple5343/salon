package integration

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	db *sql.DB
	m  *migrate.Migrate
}

func NewMigrator(dsn string, migrationDir string) (*Migrator, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	fsrc, err := (&file.File{}).Open(fmt.Sprintf("file://%s", migrationDir))
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("file", fsrc, "postgres", driver)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		db: db,
		m:  m,
	}, nil
}

func (m *Migrator) Up() error {
	return m.m.Up()
}

func (m *Migrator) Down() error {
	return m.m.Down()
}

func (m *Migrator) Close() error {
	return m.db.Close()
}
