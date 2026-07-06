package tasks

import (
	"context"
	"errors"
	"testing"
)

// fakeTaskRepo is an in-memory Repository for the tasks service/handler tests.
type fakeTaskRepo struct {
	createID  string
	createErr error
	tasks     []Task
	task      *Task
	err       error

	created     *CreateTaskParam
	gotUserID   string
	gotTaskID   string
	gotComplete bool
}

func (f *fakeTaskRepo) CreateTask(ctx context.Context, params CreateTaskParam) (string, error) {
	f.created = &params
	if f.createErr != nil {
		return "", f.createErr
	}
	return f.createID, nil
}

func (f *fakeTaskRepo) GetTasks(ctx context.Context, userId string) ([]Task, error) {
	f.gotUserID = userId
	return f.tasks, f.err
}

func (f *fakeTaskRepo) DeleteTask(ctx context.Context, userId, taskId string) error {
	f.gotUserID = userId
	f.gotTaskID = taskId
	return f.err
}

func (f *fakeTaskRepo) GetTask(ctx context.Context, userId, taskId string) (*Task, error) {
	f.gotUserID = userId
	f.gotTaskID = taskId
	return f.task, f.err
}

func (f *fakeTaskRepo) UpdateTask(ctx context.Context, userId, taskId string, isCompleted bool) (*Task, error) {
	f.gotUserID = userId
	f.gotTaskID = taskId
	f.gotComplete = isCompleted
	return f.task, f.err
}

func TestServiceCreateTask(t *testing.T) {
	repo := &fakeTaskRepo{createID: "task-id"}
	svc := newService(repo)

	id, err := svc.CreateTask(context.Background(), CreateTaskRequest{
		Title:       "write tests",
		Description: "cover the modules",
	}, "user-1")
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if id != "task-id" {
		t.Fatalf("id = %q, want task-id", id)
	}
	if repo.created == nil {
		t.Fatal("repo CreateTask not called")
	}
	if repo.created.UserId != "user-1" {
		t.Fatalf("user id = %q, want user-1", repo.created.UserId)
	}
	if repo.created.IsCompleted {
		t.Fatal("new task should not be completed")
	}
	if repo.created.Title != "write tests" || repo.created.Description != "cover the modules" {
		t.Fatalf("fields not mapped: %+v", repo.created)
	}
	if repo.created.ID == "" {
		t.Fatal("expected a generated task id")
	}
}

func TestServiceCreateTask_PropagatesError(t *testing.T) {
	svc := newService(&fakeTaskRepo{createErr: errors.New("boom")})
	if _, err := svc.CreateTask(context.Background(), CreateTaskRequest{}, "user-1"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestServiceDeleteTask_ArgOrder(t *testing.T) {
	repo := &fakeTaskRepo{}
	svc := newService(repo)

	if err := svc.DeleteTask(context.Background(), "user-1", "task-9"); err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
	if repo.gotUserID != "user-1" || repo.gotTaskID != "task-9" {
		t.Fatalf("repo got (userID=%q, taskID=%q), want (user-1, task-9)", repo.gotUserID, repo.gotTaskID)
	}
}
