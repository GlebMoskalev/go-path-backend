package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskMeta struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Order       int    `yaml:"order"`
}

type TaskFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Order       int    `yaml:"order"`
	Difficulty  string `yaml:"difficulty"`
}

type TaskChapter struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Tasks       []Task `json:"tasks"`
}

type Task struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Template    string `json:"template,omitempty"`
	Difficulty  string `json:"difficulty"`
	Order       int    `json:"order"`
	ChapterSlug string `json:"chapter_slug,omitempty"`
}

type SubmitResult struct {
	Passed bool         `json:"passed"`
	Tests  []TestResult `json:"tests"`
	Error  string       `json:"error,omitempty"`
}

type TestResult struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
	Output string `json:"output,omitempty"`
}

type Submission struct {
	ID          uuid.UUID    `json:"id"`
	UserID      uuid.UUID    `json:"user_id"`
	ChapterSlug string       `json:"chapter_slug"`
	TaskSlug    string       `json:"task_slug"`
	Code        string       `json:"code"`
	Passed      bool         `json:"passed"`
	Result      SubmitResult `json:"result"`
	CreatedAt   time.Time    `json:"created_at"`
}

type SolvedTask struct {
	ChapterSlug string `json:"chapter_slug"`
	TaskSlug    string `json:"task_slug"`
}
