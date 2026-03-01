package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	SubmissionNotFound = errors.New("submission not found")
)

type SubmissionRepository interface {
	Create(ctx context.Context, s *model.Submission) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Submission, error)
	ListByUserAndTask(ctx context.Context, userID uuid.UUID, chapterSlug, taskSlug string) ([]model.Submission, error)
	GetSolvedTasks(ctx context.Context, userID uuid.UUID) ([]model.SolvedTask, error)
	HasSolved(ctx context.Context, userID uuid.UUID, chapterSlug, taskSlug string) (bool, error)
}

type submissionRepository struct {
	db *pgxpool.Pool
}

func NewSubmissionRepository(db *pgxpool.Pool) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Create(ctx context.Context, s *model.Submission) error {
	result, err := json.Marshal(s.Result)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO submissions (user_id, chapter_slug, task_slug, code, passed, result)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query, s.UserID, s.ChapterSlug, s.TaskSlug, s.Code, s.Passed, result).
		Scan(&s.ID, &s.CreatedAt)
}

func (r *submissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Submission, error) {
	query := `
	SELECT id, user_id, chapter_slug, task_slug, code, passed, result, created_at
	FROM submissions
	WHERE id = $1
	`

	s := &model.Submission{}
	var resultJSON []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.ChapterSlug, &s.TaskSlug,
		&s.Code, &s.Passed, &resultJSON, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, SubmissionNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(resultJSON, &s.Result); err != nil {
		return nil, err
	}

	return s, nil
}

func (r *submissionRepository) ListByUserAndTask(ctx context.Context, userID uuid.UUID, chapterSlug, taskSlug string) ([]model.Submission, error) {
	query := `
	SELECT id, user_id, chapter_slug, task_slug, code, passed, result, created_at
	FROM submissions
	WHERE user_id = $1  AND chapter_slug = $2 AND task_slug = $3
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID, chapterSlug, taskSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []model.Submission
	for rows.Next() {
		var s model.Submission
		var resultJSON []byte

		if err := rows.Scan(
			&s.ID, &s.UserID, &s.ChapterSlug, &s.TaskSlug,
			&s.Code, &s.Passed, &resultJSON, &s.CreatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(resultJSON, &s.Result); err != nil {
			return nil, err
		}

		submissions = append(submissions, s)
	}

	return submissions, rows.Err()
}

func (r *submissionRepository) GetSolvedTasks(ctx context.Context, userID uuid.UUID) ([]model.SolvedTask, error) {
	query := `
	SELECT DISTINCT chapter_slug, task_slug
	FROM submissions
	WHERE user_id = $1 AND passed = TRUE
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var solvedTasks []model.SolvedTask
	for rows.Next() {
		var s model.SolvedTask

		if err := rows.Scan(&s.ChapterSlug, &s.TaskSlug); err != nil {
			return nil, err
		}

		solvedTasks = append(solvedTasks, s)
	}

	return solvedTasks, rows.Err()
}

func (r *submissionRepository) HasSolved(ctx context.Context, userID uuid.UUID, chapterSlug, taskSlug string) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1 FROM submissions
		WHERE user_id = $1 AND chapter_slug = $2 AND task_slug = $3 AND passed = TRUE
	)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, chapterSlug, taskSlug).Scan(&exists)
	return exists, err
}
