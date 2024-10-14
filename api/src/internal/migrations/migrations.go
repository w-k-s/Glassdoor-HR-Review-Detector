package migrations

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Exec(migrationsDirectory string, db *sql.DB) error {
	if len(migrationsDirectory) == 0 {
		return fmt.Errorf("invalid migrations directory: '%s'. Must be an absolute path", migrationsDirectory)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("Failed to create migrations driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDirectory,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Failed to load migrations directory '%s': %w", migrationsDirectory, err)
	}
	return m.Up()
}
