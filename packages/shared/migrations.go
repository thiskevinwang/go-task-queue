package shared

import (
	"shared/log"

	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(logger log.Logger) error {
	db := NewDB().DB
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	if err != nil {

		// This is simple heuristic for incorrectly named migration file
		if strings.Contains(err.Error(), "file does not exist") {
			logger.Error("Possible incorrectly named migration file. Make sure the file is named like \"foobar.up.sql\"", "err", err.Error())
			return err
		} else if err == migrate.ErrNoChange {
			logger.Info("Skipping migrations", "err", err)
		} else {
			logger.Error("Migrations failed", "err", err)
			return err
		}
	}
	return nil
}
