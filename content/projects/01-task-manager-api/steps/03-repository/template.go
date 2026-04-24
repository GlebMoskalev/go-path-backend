package repository

// Реализуйте TaskRepository — слой работы с PostgreSQL.
//
// Импорты:
//   import (
//       "context"
//       "database/sql"
//       "fmt"
//       "taskmanager/model"
//   )
//
// 1. Определите структуру:
//    type TaskRepository struct {
//        db *sql.DB
//    }
//
// 2. Конструктор:
//    func New(db *sql.DB) *TaskRepository
//
// 3. Методы (все принимают context.Context первым аргументом):
//
//    Create(ctx, task model.Task) (model.Task, error)
//    — INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3)
//    — RETURNING id, created_at, updated_at  (используй QueryRowContext + Scan)
//
//    GetByID(ctx, id int64) (model.Task, error)
//    — SELECT ... FROM tasks WHERE id = $1
//    — при sql.ErrNoRows: return model.Task{}, fmt.Errorf("task %d: not found", id)
//
//    List(ctx) ([]model.Task, error)
//    — SELECT ... FROM tasks ORDER BY created_at DESC
//    — инициализируй: tasks := make([]model.Task, 0)
//    — проверяй rows.Err() после for rows.Next()
//
//    Update(ctx, task model.Task) (model.Task, error)
//    — UPDATE tasks SET title=$1, description=$2, status=$3, updated_at=NOW() WHERE id=$4
//    — RETURNING updated_at
//
//    Delete(ctx, id int64) error
//    — DELETE FROM tasks WHERE id = $1
//    — проверяй result.RowsAffected() — если 0, верни fmt.Errorf("task %d: not found", id)
