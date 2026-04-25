package repository

import (
	"context"
	"database/sql"
	"fmt"
	"taskmanager/model"
)

type TaskRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task model.Task) (model.Task, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO tasks (title, description, status)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		task.Title, task.Description, task.Status,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	return task, err
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (model.Task, error) {
	var task model.Task
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, status, created_at, updated_at
		 FROM tasks WHERE id = $1`,
		id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return model.Task{}, fmt.Errorf("task %d: not found", id)
	}
	return task, err
}

func (r *TaskRepository) List(ctx context.Context) ([]model.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, description, status, created_at, updated_at
		 FROM tasks ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description,
			&task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) Update(ctx context.Context, task model.Task) (model.Task, error) {
	err := r.db.QueryRowContext(ctx,
		`UPDATE tasks SET title=$1, description=$2, status=$3, updated_at=NOW()
		 WHERE id=$4
		 RETURNING updated_at`,
		task.Title, task.Description, task.Status, task.ID,
	).Scan(&task.UpdatedAt)
	if err == sql.ErrNoRows {
		return model.Task{}, fmt.Errorf("task %d: not found", task.ID)
	}
	return task, err
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("task %d: not found", id)
	}
	return nil
}
