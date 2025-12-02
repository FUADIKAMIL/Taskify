package repository

import (
    "context"
    "database/sql"
    "errors"

    "github.com/FUADIKAMIL/taskify/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
    Create        func(ctx context.Context, u model.User) (model.User, error)
    GetByUsername func(ctx context.Context, username string) (model.User, error)
    GetByID       func(ctx context.Context, id int64) (model.User, error)
}

func NewUserRepo(db *sql.DB) *UserRepo {

    return &UserRepo{

        Create: func(ctx context.Context, u model.User) (model.User, error) {
            var id int64
            err := db.QueryRowContext(
                ctx,
                `INSERT INTO users (username, password, created_at)
                 VALUES ($1,$2,now()) RETURNING id`,
                u.Username, u.Password,
            ).Scan(&id)

            if err != nil {
                return model.User{}, err
            }

            u.ID = id
            return u, nil
        },

        GetByUsername: func(ctx context.Context, username string) (model.User, error) {
            var u model.User
            err := db.QueryRowContext(
                ctx,
                `SELECT id, username, password, created_at 
                 FROM users WHERE username=$1`,
                username,
            ).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)

            if err != nil {
                if err == sql.ErrNoRows {
                    return model.User{}, ErrUserNotFound
                }
                return model.User{}, err
            }
            return u, nil
        },

        GetByID: func(ctx context.Context, id int64) (model.User, error) {
            var u model.User
            err := db.QueryRowContext(
                ctx,
                `SELECT id, username, password, created_at 
                 FROM users WHERE id=$1`,
                id,
            ).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)

            if err != nil {
                if err == sql.ErrNoRows {
                    return model.User{}, ErrUserNotFound
                }
                return model.User{}, err
            }
            return u, nil
        },
    }
}
