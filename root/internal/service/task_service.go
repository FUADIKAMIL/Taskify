package service

import (
    "context"
    "time"

    "github.com/yourname/taskify/internal/model"
    "github.com/yourname/taskify/internal/repository"
)

type TaskService struct{ repo *repository.TaskRepo }

func NewTaskService(r *repository.TaskRepo) *TaskService { return &TaskService{repo: r} }

func computeStatus(t model.Task, now time.Time) string {
    if t.Completed {
        return "completed"
    }
    if t.Deadline != nil && t.Deadline.Before(now) {
        return "overdue"
    }
    return "pending"
}

func (s *TaskService) CreateTask(ctx context.Context, t model.Task) (model.Task, error) {
    created, err := s.repo.Create(ctx, t)
    if err != nil {
        return model.Task{}, err
    }
    created.Status = computeStatus(created, time.Now())
    return created, nil
}

func (s *TaskService) GetAllTasks(ctx context.Context, userID int64) ([]model.Task, error) {
    list, err := s.repo.GetAllByUser(ctx, userID)
    if err != nil {
        return nil, err
    }
    now := time.Now()
    for i := range list {
        list[i].Status = computeStatus(list[i], now)
    }
    return list, nil
}

func (s *TaskService) GetTask(ctx context.Context, id, userID int64) (model.Task, error) {
    t, err := s.repo.GetByIDAndUser(ctx, id, userID)
    if err != nil {
        return model.Task{}, err
    }
    t.Status = computeStatus(t, time.Now())
    return t, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, t model.Task) (model.Task, error) {
    updated, err := s.repo.Update(ctx, t)
    if err != nil {
        return model.Task{}, err
    }
    updated.Status = computeStatus(updated, time.Now())
    return updated, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id, userID int64) error {
    return s.repo.Delete(ctx, id, userID)
}
