package model

type QuizQuestion struct {
	ID          string   `json:"id" yaml:"-"`
	Question    string   `json:"question" yaml:"question"`
	Options     []string `json:"options" yaml:"options"`
	Answer      int      `json:"-" yaml:"answer"`
	Explanation string   `json:"-" yaml:"explanation"`
	ChapterSlug string   `json:"chapter_slug" yaml:"-"`
}

type QuizChapterInfo struct {
	Slug          string `json:"slug"`
	Title         string `json:"title"`
	QuestionCount int    `json:"question_count"`
}

type QuizAnswerResponse struct {
	Correct       bool   `json:"correct"`
	CorrectAnswer int    `json:"correct_answer"`
	Explanation   string `json:"explanation"`
}
