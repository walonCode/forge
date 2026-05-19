package database

import (
	"database/sql"
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

	return &Database{
		Client: db,
	}, nil
}
