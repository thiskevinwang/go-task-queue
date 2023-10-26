package shared

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/hashicorp/go-hclog"
)

func RunMigrations(log hclog.Logger) {
	db := NewDB().DB
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	err = m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Info("Skipping migrations", "err", err)
		} else {
			log.Error("Migrations failed", "err", err)
			panic(err)
		}
	}
}
