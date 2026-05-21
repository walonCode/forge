package tasks

import (
	"context"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Service struct {
	repo Respository
}

func newService(repo Respository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTask(ctx context.Context, params CreateTaskRequest, userId string) (string, error) {
	id, err := gonanoid.New(30)
	if err != nil {
		return "", err
	}

	newTask := CreateTaskParam{
		ID:          id,
		IsCompleted: false,
		Title:       params.Title,
		Description: params.Description,
		UserId:      userId,
	}

	data, err := s.repo.CreateTask(ctx, newTask)
	if err != nil {
		return "", err
	}

	return data, nil
}

func (s *Service) GetTasks(ctx context.Context, userId string) (*[]Task, error) {
	tasks, err := s.repo.GetTasks(ctx, userId)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) DeleteTask(ctx context.Context, userId, taskId string) error {
	if err := s.repo.DeleteTask(ctx, userId, taskId); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTask(ctx context.Context, userId, taskId string) (*Task, error) {
	task, err := s.repo.GetTask(ctx, userId, taskId)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, userId, taskId string, isCompleted bool) (*Task, error) {
	task, err := s.repo.UpdateTask(ctx, userId, taskId, isCompleted)
	if err != nil {
		return nil, err
	}

	return task, nil
}
