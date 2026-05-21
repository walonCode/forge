package tasks

import (
	"context"
	"database/sql"
)


type Respository interface {
	CreateTask(ctx context.Context, params CreateTaskParam)(string, error)
	GetTasks(ctx context.Context, userId string)(*[]Task, error)
	DeleteTask(ctx context.Context, userId, taskId string)(error)
}

type sqlRepository struct {
	db *sql.DB
}

func newRepository(db *sql.DB)Respository{
	return &sqlRepository{
		db: db,
	}
}

func (r *sqlRepository)CreateTask(ctx context.Context, params CreateTaskParam)(string, error){
	var taskId string
	sql := "INSERT INTO tasks (id, title, userId, description isCompleted) VALUES ($1,$2,$3,$4,$5) RETURNING id"
	
	//adding it to the db
	if err := r.db.QueryRowContext(ctx, sql,params.ID, params.Title, params.UserId, params.Description, params.IsCompleted ).Scan(&taskId); err != nil {
		return "", err
	}

	return taskId, nil
}

func (r *sqlRepository)GetTasks(ctx context.Context, userId string)(*[]Task, error){
	var tasks []Task
	sql := "SELECT * FROM tasks WHERE userId = $1"

	rows, err := r.db.QueryContext(ctx, sql, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.IsCompleted, &task.UserId, &task.UpdatedAt, &task.CreatedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	defer rows.Close()

	return &tasks, nil
}

func (r *sqlRepository)DeleteTask(ctx context.Context, userId, taskId string)(error){
	sql := "DELETE FROM tasks WHERE userId = $1 AND taskId = $2"

	if _, err := r.db.ExecContext(ctx, sql, userId, taskId); err != nil {
		return err
	}

	return nil
}