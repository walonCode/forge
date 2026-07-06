package tasks

import (
	"context"
	"database/sql"
)

type Repository interface {
	CreateTask(ctx context.Context, params CreateTaskParam) (string, error)
	GetTasks(ctx context.Context, userId string) ([]Task, error)
	DeleteTask(ctx context.Context, userId, taskId string) error
	GetTask(ctx context.Context, userId, taskId string) (*Task, error)
	UpdateTask(ctx context.Context, userId, taskId string, isCompleted bool) (*Task, error)
}

type sqlRepository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) Repository {
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository) CreateTask(ctx context.Context, params CreateTaskParam) (string, error) {
	var taskId string
	query := "INSERT INTO tasks (id, title, userId, description, is_completed) VALUES ($1,$2,$3,$4,$5) RETURNING id"

	//adding it to the db
	if err := r.db.QueryRowContext(ctx, query, params.ID, params.Title, params.UserId, params.Description, params.IsCompleted).Scan(&taskId); err != nil {
		return "", err
	}

	return taskId, nil
}

func (r *sqlRepository) GetTasks(ctx context.Context, userId string) ([]Task, error) {
	var tasks []Task
	query := "SELECT id,title,description, is_completed, userId, updated_at, created_at FROM tasks WHERE userId = $1"

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.UserId, &task.UpdatedAt, &task.CreatedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *sqlRepository) DeleteTask(ctx context.Context, userId, taskId string) error {
	query := "DELETE FROM tasks WHERE userId = $1 AND id = $2"

	if _, err := r.db.ExecContext(ctx, query, userId, taskId); err != nil {
		return err
	}

	return nil
}

func (r *sqlRepository) GetTask(ctx context.Context, userId, taskId string) (*Task, error) {
	query := "SELECT id,title,description, is_completed, updated_at, userId, created_at  FROM tasks WHERE userId = $1 AND id = $2"

	var task Task
	if err := r.db.QueryRowContext(ctx, query, userId, taskId).Scan(
		&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.UpdatedAt, &task.UserId, &task.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *sqlRepository) UpdateTask(ctx context.Context, userId, taskId string, isCompleted bool) (*Task, error) {
	query := "UPDATE tasks SET is_completed = $1, updated_at = now() WHERE userId = $2 AND id = $3 RETURNING id, title, description, is_completed, updated_at, userId, created_at"

	var task Task
	if err := r.db.QueryRowContext(ctx, query, isCompleted, userId, taskId).Scan(
		&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.UpdatedAt, &task.UserId, &task.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &task, nil
}
