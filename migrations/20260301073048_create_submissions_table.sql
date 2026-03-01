-- +goose Up
-- +goose StatementBegin
CREATE TABLE submissions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chapter_slug VARCHAR(100) NOT NULL,
    task_slug    VARCHAR(100) NOT NULL,
    code       TEXT NOT NULL,
    passed     BOOLEAN NOT NULL DEFAULT FALSE,
    result     JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_submissions_user_task ON submissions(user_id, chapter_slug, task_slug)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS submissions;
-- +goose StatementEnd
