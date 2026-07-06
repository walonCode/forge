package auth

import (
	"context"
	"database/sql"
)

type Repository interface {
	CreateUser(ctx context.Context, params CreateDbUser) (string, error)
	FindUserByUsername(ctx context.Context, username string) (*DbUser, error)
}

type sqlRepository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) Repository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) CreateUser(ctx context.Context, params CreateDbUser) (string, error) {
	var Id string
	query := "INSERT INTO users (id, name, username, password) VALUES ($1, $2, $3, $4) RETURNING id"

	if err := r.db.QueryRowContext(ctx, query, params.ID, params.Name, params.Username, params.Password).Scan(&Id); err != nil {
		return "", err
	}

	return Id, nil
}

func (r *sqlRepository) FindUserByUsername(ctx context.Context, username string) (*DbUser, error) {
	var user DbUser
	query := "SELECT password, username, id FROM users WHERE username = $1"

	if err := r.db.QueryRowContext(ctx, query, username).Scan(&user.Password, &user.Username, &user.ID); err != nil {
		return nil, err
	}

	return &user, nil
}
