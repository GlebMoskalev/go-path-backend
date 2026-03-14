-- +goose Up
-- +goose StatementBegin
CREATE TABLE theories_progress(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chapter_slug VARCHAR(100) NOT NULL,
    lesson_slug VARCHAR(100) NOT NULL,
    completed_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(user_id, chapter_slug, lesson_slug)
);

CREATE INDEX idx_lesson_progress_user ON theories_progress(user_id);
CREATE INDEX idx_lesson_progress_user_chapter ON theories_progress(user_id, chapter_slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lesson_progress;
-- +goose StatementEnd
