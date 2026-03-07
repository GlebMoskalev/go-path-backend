package model

type ProjectMeta struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Order       int    `yaml:"order"`
}

type StepFrontmatter struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Order       int      `yaml:"order"`
	Difficulty  string   `yaml:"difficulty"`
	File        string   `yaml:"file"`
	Hints       []string `yaml:"hints"`
}

type Project struct {
	Slug        string        `json:"slug"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	Order       int           `json:"order"`
	Steps       []ProjectStep `json:"steps"`
	SolvedCount int           `json:"solved_count"`
}

type ProjectStep struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Template    string   `json:"template,omitempty"`
	Difficulty  string   `json:"difficulty"`
	Hints       []string `json:"hints,omitempty"`
	File        string   `json:"file,omitempty"`
	Order       int      `json:"order"`
	ProjectSlug string   `json:"project_slug,omitempty"`
	Solved      *bool    `json:"solved,omitempty"`
}
