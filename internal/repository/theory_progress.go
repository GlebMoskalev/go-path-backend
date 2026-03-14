package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TheoryProgressRepository interface {
	MarkCompleted(ctx context.Context, userID uuid.UUID, chapterSlug, lessonSlug string) error
	GetCompletedTheories(ctx context.Context, userID uuid.UUID) (map[string]map[string]bool, error)
	IsCompleted(ctx context.Context, userID uuid.UUID, chapterSlug, lessonSlug string) (bool, error)
}

type theoryProgressRepository struct {
	db *pgxpool.Pool
}

func NewTheoryProgressRepository(db *pgxpool.Pool) TheoryProgressRepository {
	return &theoryProgressRepository{db: db}
}

func (r *theoryProgressRepository) MarkCompleted(ctx context.Context, userID uuid.UUID, chapterSlug, lessonSlug string) error {
	query := `
	INSERT INTO theories_progress (user_id, chapter_slug, lesson_slug)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id, chapter_slug, lesson_slug) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, userID, chapterSlug, lessonSlug)
	return err
}

func (r *theoryProgressRepository) GetCompletedTheories(ctx context.Context, userID uuid.UUID) (map[string]map[string]bool, error) {
	query := `
	SELECT chapter_slug, lesson_slug
	FROM theories_progress
	WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	completed := make(map[string]map[string]bool)

	for rows.Next() {
		var chapterSlug, lessonSlug string
		if err := rows.Scan(&chapterSlug, &lessonSlug); err != nil {
			return nil, err
		}

		if completed[chapterSlug] == nil {
			completed[chapterSlug] = make(map[string]bool)
		}
		completed[chapterSlug][lessonSlug] = true
	}

	return completed, rows.Err()
}

func (r *theoryProgressRepository) IsCompleted(ctx context.Context, userID uuid.UUID, chapterSlug, lessonSlug string) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1 FROM theories_progress
		WHERE user_id = $1 AND chapter_slug = $2 AND lesson_slug = $3
	)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, chapterSlug, lessonSlug).Scan(&exists)
	return exists, err
}
