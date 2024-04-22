package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dsn   = "phrase.db?cache=shared&mode=r"
	MaxID = 467014991
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
