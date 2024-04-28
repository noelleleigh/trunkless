package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const (
	dburl = "postgresql://vilmibm/postgres?host=/home/vilmibm/src/trunkless/pgdata/sockets"
	//	MaxID = 467014991
	MaxID = 345507789
)

func Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "")
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	return conn, nil
}

// TODO func for getting ID ranges for each corpus in phrases
