package db

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// TODO unused; i think i failed to make this work. should i just rely on the env
	// vars like i've been doing?
	dburl = "postgresql://vilmibm/postgres?host=/home/vilmibm/src/trunkless/pgdata/sockets"
)

func StrToID(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))[0:6]
}

func Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "")
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	return conn, nil
}

func Pool() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), "")
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	return pool, nil
}

// TODO func for getting ID ranges for each corpus in phrases
