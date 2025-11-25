package model

import "time"

type Task struct {
    ID        int64      `json:"id"`
    UserID    int64      `json:"user_id"`
    Title     string     `json:"title"`
    Content   string     `json:"content,omitempty"`
    Deadline  *time.Time `json:"deadline,omitempty"`
    Completed bool       `json:"completed"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    Status    string     `json:"status,omitempty"`
}
