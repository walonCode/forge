package users

import (
	"context"
	"database/sql"
)

type Repository interface {
	FindProfileByID(ctx context.Context, id string) (*UserProfile, error)
	GetPasswordHash(ctx context.Context, id string) (string, error)
	ExistsByUsernameExcept(ctx context.Context, username, excludeID string) (bool, error)
	UpdateProfile(ctx context.Context, id, name, username string) (*UserProfile, error)
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	DeleteUser(ctx context.Context, id string) error
}

type sqlRepository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) Repository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) FindProfileByID(ctx context.Context, id string) (*UserProfile, error) {
	query := "SELECT id, name, username, created_at FROM users WHERE id = $1"

	var user UserProfile
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Username, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *sqlRepository) GetPasswordHash(ctx context.Context, id string) (string, error) {
	query := "SELECT password FROM users WHERE id = $1"

	var hash string
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&hash); err != nil {
		return "", err
	}

	return hash, nil
}

func (r *sqlRepository) ExistsByUsernameExcept(ctx context.Context, username, excludeID string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 AND id <> $2)"

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, username, excludeID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *sqlRepository) UpdateProfile(ctx context.Context, id, name, username string) (*UserProfile, error) {
	query := "UPDATE users SET name = $1, username = $2, updated_at = now() WHERE id = $3 RETURNING id, name, username, created_at"

	var user UserProfile
	if err := r.db.QueryRowContext(ctx, query, name, username, id).Scan(&user.ID, &user.Name, &user.Username, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *sqlRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	query := "UPDATE users SET password = $1, updated_at = now() WHERE id = $2"

	if _, err := r.db.ExecContext(ctx, query, passwordHash, id); err != nil {
		return err
	}

	return nil
}

func (r *sqlRepository) DeleteUser(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
