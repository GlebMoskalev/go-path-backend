package model

type UserStats struct {
	Theory   TheoryStats   `json:"theory"`
	Tasks    TasksStats    `json:"tasks"`
	Projects ProjectsStats `json:"projects"`
}

type TheoryStats struct {
	TotalLessons     int                  `json:"total_lessons"`
	CompletedLessons int                  `json:"completed_lessons"`
	Chapters         []TheoryChapterStats `json:"chapters"`
}

type TheoryChapterStats struct {
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Total     int    `json:"total"`
	Completed int    `json:"completed"`
}

type TasksStats struct {
	TotalTasks  int                `json:"total_tasks"`
	SolvedTasks int                `json:"solved_tasks"`
	Chapters    []TaskChapterStats `json:"chapters"`
}

type TaskChapterStats struct {
	Slug   string `json:"slug"`
	Title  string `json:"title"`
	Total  int    `json:"total"`
	Solved int    `json:"solved"`
}

type ProjectsStats struct {
	TotalSteps  int                `json:"total_steps"`
	SolvedSteps int                `json:"solved_steps"`
	Projects    []ProjectStatsItem `json:"projects"`
}

type ProjectStatsItem struct {
	Slug   string `json:"slug"`
	Title  string `json:"title"`
	Total  int    `json:"total"`
	Solved int    `json:"solved"`
}
