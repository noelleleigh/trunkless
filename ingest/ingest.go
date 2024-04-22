package ingest

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"

	_ "github.com/mattn/go-sqlite3"
)

const dsn = "phrase.db?_journal=OFF"

func createSource(db *sql.DB, sourceName string) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO sources (name) VALUES (?) ON CONFLICT DO NOTHING RETURNING id")
	if err != nil {
		return -1, err
	}

	result, err := stmt.Exec(sourceName)
	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	return result.LastInsertId()
}

func Ingest(sourceName string, ins io.Reader) error {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	s := bufio.NewScanner(ins)

	sourceID, err := createSource(db, sourceName)
	if err != nil {
		return fmt.Errorf("could not make source: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO phrases (sourceid, phrase) VALUES (?, ?) ON CONFLICT DO NOTHING")
	defer stmt.Close()
	for s.Scan() {
		phrase := s.Text()
		if err != nil {
			return err
		}

		if _, err = stmt.Exec(sourceID, phrase); err != nil {
			return fmt.Errorf("could not insert phrase '%s' for source '%d': %w", phrase, sourceID, err)
		}
	}

	return tx.Commit()
}
