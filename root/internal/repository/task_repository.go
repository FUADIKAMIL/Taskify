package repository

import (
    "context"
    "database/sql"
    "time"

    "github.com/FUADIKAMIL/taskify/internal/model"
)

type TaskRepo struct {
    Create         func(ctx context.Context, t model.Task) (model.Task, error)
    GetAllByUser   func(ctx context.Context, userID int64) ([]model.Task, error)
    GetByIDAndUser func(ctx context.Context, id int64, userID int64) (model.Task, error)
    Update         func(ctx context.Context, t model.Task) (model.Task, error)
    Delete         func(ctx context.Context, id int64, userID int64) error
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
    return &TaskRepo{

        Create: func(ctx context.Context, t model.Task) (model.Task, error) {
            now := time.Now()
            var id int64

            err := db.QueryRowContext(ctx,
                `INSERT INTO tasks (user_id, title, content, deadline, completed, created_at, updated_at)
                 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
                t.UserID, t.Title, t.Content, t.Deadline, t.Completed, now, now,
            ).Scan(&id)

            if err != nil {
                return model.Task{}, err
            }

            t.ID = id
            t.CreatedAt = now
            t.UpdatedAt = now
            return t, nil
        },

        GetAllByUser: func(ctx context.Context, userID int64) ([]model.Task, error) {
            rows, err := db.QueryContext(ctx,
                `SELECT id, user_id, title, content, deadline, completed, created_at, updated_at
                 FROM tasks
                 WHERE user_id=$1
                 ORDER BY deadline ASC NULLS LAST`,
                userID,
            )
            if err != nil {
                return nil, err
            }
            defer rows.Close()

            var out []model.Task

            for rows.Next() {
                var t model.Task
                var dt sql.NullTime

                if err := rows.Scan(
                    &t.ID, &t.UserID, &t.Title, &t.Content,
                    &dt, &t.Completed, &t.CreatedAt, &t.UpdatedAt,
                ); err != nil {
                    return nil, err
                }

                if dt.Valid {
                    d := dt.Time
                    t.Deadline = &d
                }

                out = append(out, t)
            }

            return out, nil
        },

        GetByIDAndUser: func(ctx context.Context, id int64, userID int64) (model.Task, error) {
            var t model.Task
            var dt sql.NullTime

            err := db.QueryRowContext(ctx,
                `SELECT id, user_id, title, content, deadline, completed, created_at, updated_at
                 FROM tasks
                 WHERE id=$1 AND user_id=$2`,
                id, userID,
            ).Scan(&t.ID, &t.UserID, &t.Title, &t.Content,
                &dt, &t.Completed, &t.CreatedAt, &t.UpdatedAt)

            if err != nil {
                return model.Task{}, err
            }

            if dt.Valid {
                d := dt.Time
                t.Deadline = &d
            }

            return t, nil
        },

        Update: func(ctx context.Context, t model.Task) (model.Task, error) {
            now := time.Now()

            _, err := db.ExecContext(ctx,
                `UPDATE tasks
                 SET title=$1, content=$2, deadline=$3, completed=$4, updated_at=$5
                 WHERE id=$6 AND user_id=$7`,
                t.Title, t.Content, t.Deadline, t.Completed, now, t.ID, t.UserID,
            )

            if err != nil {
                return model.Task{}, err
            }

            t.UpdatedAt = now
            return t, nil
        },

        Delete: func(ctx context.Context, id int64, userID int64) error {
            _, err := db.ExecContext(ctx,
                `DELETE FROM tasks WHERE id=$1 AND user_id=$2`,
                id, userID,
            )
            return err
        },
    }
}
