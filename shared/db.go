package shared

import (
	"database/sql"
	"os"
)

type DB struct {
	*sql.DB
}

func NewDB() *DB {
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}

	return &DB{
		db,
	}
}
