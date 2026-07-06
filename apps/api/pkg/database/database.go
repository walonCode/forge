package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Client *sql.DB
}

func (d *Database) Close() error {
	if d.Client != nil {
		return d.Client.Close()
	}

	return nil
}

func Connect(database_url string) (*Database, error) {
	db, err := sql.Open("pgx", database_url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// verify the connection is actually reachable so misconfiguration fails fast
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{
		Client: db,
	}, nil
}
