package service

import (
    "context"
    "time"

    "github.com/FUADIKAMIL/taskify/internal/model"
    "github.com/FUADIKAMIL/taskify/internal/repository"
)

func computeStatus(t model.Task, now time.Time) string {
    if t.Completed {
        return "completed"
    }
    if t.Deadline != nil && t.Deadline.Before(now) {
        return "overdue"
    }
    return "pending"
}

type TaskService struct {
    CreateTask func(ctx context.Context, t model.Task) (model.Task, error)
    GetAllTasks func(ctx context.Context, userID int64) ([]model.Task, error)
    GetTask func(ctx context.Context, id, userID int64) (model.Task, error)
    UpdateTask func(ctx context.Context, t model.Task) (model.Task, error)
    DeleteTask func(ctx context.Context, id, userID int64) error
}

func NewTaskService(repo *repository.TaskRepo) *TaskService {

    return &TaskService{

        CreateTask: func(ctx context.Context, t model.Task) (model.Task, error) {
            created, err := repo.Create(ctx, t)
            if err != nil {
                return model.Task{}, err
            }
            created.Status = computeStatus(created, time.Now())
            return created, nil
        },

        GetAllTasks: func(ctx context.Context, userID int64) ([]model.Task, error) {
            list, err := repo.GetAllByUser(ctx, userID)
            if err != nil {
                return nil, err
            }
            now := time.Now()
            for i := range list {
                list[i].Status = computeStatus(list[i], now)
            }
            return list, nil
        },

        GetTask: func(ctx context.Context, id, userID int64) (model.Task, error) {
            t, err := repo.GetByIDAndUser(ctx, id, userID)
            if err != nil {
                return model.Task{}, err
            }
            t.Status = computeStatus(t, time.Now())
            return t, nil
        },

        UpdateTask: func(ctx context.Context, t model.Task) (model.Task, error) {
            updated, err := repo.Update(ctx, t)
            if err != nil {
                return model.Task{}, err
            }
            updated.Status = computeStatus(updated, time.Now())
            return updated, nil
        },

        DeleteTask: func(ctx context.Context, id, userID int64) error {
            return repo.Delete(ctx, id, userID)
        },
    }
}
