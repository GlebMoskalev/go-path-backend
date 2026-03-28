package model

type Chapter struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Order       int      `json:"order"`
	Lessons     []Lesson `json:"lessons"`
	// Progress nil = не авторизован, иначе статистика
	Progress *ChapterProgress `json:"progress,omitempty"`
}

type ChapterProgress struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
}

type Lesson struct {
	Slug         string `json:"slug"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Order        int    `json:"order"`
	ChapterSlug  string `json:"chapter_slug,omitempty"`
	ChapterTitle string `json:"chapter_title,omitempty"`
	Content      string `json:"content,omitempty"`
	// Completed: nil = не авторизован, true/false = авторизован
	Completed *bool `json:"completed,omitempty"`
}

type ChapterMeta struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Order       int    `yaml:"order"`
}

type LessonFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Order       int    `yaml:"order"`
}
