package model

type Chapter struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Order       int      `json:"order"`
	Lessons     []Lesson `json:"lessons"`
}

type Lesson struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	ChapterSlug string `json:"chapter_slug,omitempty"`
	Content     string `json:"content,omitempty"`
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
