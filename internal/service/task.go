package service

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"github.com/GlebMoskalev/go-path-backend/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	ErrTaskChapterNotFound = errors.New("task chapter not found")
	ErrTaskNotFound        = errors.New("task not found")
)

type TaskService struct {
	log            *zap.Logger
	submissionRepo repository.SubmissionRepository
	chapters       []model.TaskChapter
	tasks          map[string]map[string]model.Task
	tests          map[string]map[string]string
}

func NewTaskService(fsys fs.FS, root string, log *zap.Logger, submissionRepo repository.SubmissionRepository) (*TaskService, error) {
	s := &TaskService{
		log:            log,
		tasks:          make(map[string]map[string]model.Task),
		tests:          make(map[string]map[string]string),
		submissionRepo: submissionRepo,
	}

	if err := s.load(fsys, root); err != nil {
		return nil, err
	}

	return s, nil
}

// ListChapters — все главы задач, задачи без description/template
func (s *TaskService) ListChapters(ctx context.Context, userID *uuid.UUID) []model.TaskChapter {
	result := make([]model.TaskChapter, len(s.chapters))
	for i, ch := range s.chapters {
		result[i] = model.TaskChapter{
			Slug:        ch.Slug,
			Title:       ch.Title,
			Description: ch.Description,
			Order:       ch.Order,
			Tasks:       stripTaskContent(ch.Tasks),
		}
	}

	if userID != nil {
		solved := s.getSolvedSet(ctx, *userID)
		for i := range result {
			count := 0
			for j := range result[i].Tasks {
				v := solved[result[i].Slug+"/"+result[i].Tasks[j].Slug]
				result[i].Tasks[j].Solved = &v
				if v {
					count++
				}
			}
			result[i].SolvedCount = count
		}
	}

	return result
}

// GetChapter — одна глава, задачи без description/template
func (s *TaskService) GetChapter(ctx context.Context, slug string, userID *uuid.UUID) (model.TaskChapter, error) {
	for _, ch := range s.chapters {
		if ch.Slug == slug {
			result := model.TaskChapter{
				Slug:        ch.Slug,
				Title:       ch.Title,
				Description: ch.Description,
				Order:       ch.Order,
				Tasks:       stripTaskContent(ch.Tasks),
			}
			if userID != nil {
				solvedCount := s.enrichWithSolved(ctx, *userID, &result)
				result.SolvedCount = solvedCount
			}
			return result, nil
		}
	}
	return model.TaskChapter{}, ErrTaskChapterNotFound
}

// GetTask — одна задача С description и template (для страницы задачи)
func (s *TaskService) GetTask(ctx context.Context, chapterSlug, taskSlug string, userID *uuid.UUID) (model.Task, error) {
	chapterTasks, ok := s.tasks[chapterSlug]
	if !ok {
		return model.Task{}, ErrTaskChapterNotFound
	}
	task, ok := chapterTasks[taskSlug]
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}
	if userID != nil {
		solved, _ := s.submissionRepo.HasSolved(ctx, *userID, chapterSlug, taskSlug)
		task.Solved = &solved
	}
	return task, nil
}

// GetTestFile — содержимое solution_test.go (только для SandboxService, НЕ для API)
func (s *TaskService) GetTestFile(chapterSlug, taskSlug string) (string, error) {
	chapterTests, ok := s.tests[chapterSlug]
	if !ok {
		return "", ErrTaskChapterNotFound
	}
	test, ok := chapterTests[taskSlug]
	if !ok {
		return "", ErrTaskNotFound
	}
	return test, nil
}

func (s *TaskService) load(fsys fs.FS, root string) error {
	chapterDirs, err := fs.ReadDir(fsys, root)
	if err != nil {
		return err
	}

	for _, dir := range chapterDirs {
		if !dir.IsDir() {
			continue
		}

		chapter, tasks, tests, err := s.loadChapter(fsys, root, dir.Name())
		if err != nil {
			s.log.Warn("skipping task chapter", zap.String("dir", dir.Name()), zap.Error(err))
			continue
		}

		chapter.Tasks = tasks
		s.chapters = append(s.chapters, chapter)

		s.tasks[chapter.Slug] = make(map[string]model.Task)
		s.tests[chapter.Slug] = make(map[string]string)
		for _, t := range tasks {
			s.tasks[chapter.Slug][t.Slug] = t
		}
		for slug, content := range tests {
			s.tests[chapter.Slug][slug] = content
		}
	}

	sort.Slice(s.chapters, func(i, j int) bool {
		return s.chapters[i].Order < s.chapters[j].Order
	})

	s.log.Info("tasks loaded",
		zap.Int("chapters", len(s.chapters)),
		zap.Int("total_tasks", s.totalTasks()),
	)

	return nil
}

func (s *TaskService) loadChapter(fsys fs.FS, root, dirName string) (model.TaskChapter, []model.Task, map[string]string, error) {
	chapterPath := filepath.Join(root, dirName)

	metaData, err := fs.ReadFile(fsys, filepath.Join(chapterPath, "meta.yaml"))
	if err != nil {
		return model.TaskChapter{}, nil, nil, err
	}

	var meta model.TaskMeta
	if err := yaml.Unmarshal(metaData, &meta); err != nil {
		return model.TaskChapter{}, nil, nil, err
	}

	chapter := model.TaskChapter{
		Slug:        dirName,
		Title:       meta.Title,
		Description: meta.Description,
		Order:       meta.Order,
	}

	taskDirs, err := fs.ReadDir(fsys, chapterPath)
	if err != nil {
		return model.TaskChapter{}, nil, nil, err
	}

	var tasks []model.Task
	tests := make(map[string]string)

	for _, td := range taskDirs {
		if !td.IsDir() {
			continue
		}

		task, testContent, err := s.loadTask(fsys, chapterPath, td.Name(), dirName)
		if err != nil {
			s.log.Warn("skipping task", zap.String("dir", td.Name()), zap.Error(err))
			continue
		}

		tasks = append(tasks, task)
		tests[task.Slug] = testContent
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Order < tasks[j].Order
	})

	return chapter, tasks, tests, nil
}

func (s *TaskService) loadTask(fsys fs.FS, chapterPath, dirName, chapterSlug string) (model.Task, string, error) {
	taskPath := filepath.Join(chapterPath, dirName)

	taskData, err := fs.ReadFile(fsys, filepath.Join(taskPath, "task.md"))
	if err != nil {
		return model.Task{}, "", err
	}

	fm, description, err := parseTaskFrontmatter(string(taskData))
	if err != nil {
		return model.Task{}, "", err
	}

	templateData, err := fs.ReadFile(fsys, filepath.Join(taskPath, "template.go"))
	if err != nil {
		return model.Task{}, "", err
	}

	testData, err := fs.ReadFile(fsys, filepath.Join(taskPath, "solution_test.go"))
	if err != nil {
		return model.Task{}, "", err
	}

	task := model.Task{
		Slug:        dirName,
		Title:       fm.Title,
		Description: description,
		Template:    string(templateData),
		Difficulty:  fm.Difficulty,
		Order:       fm.Order,
		ChapterSlug: chapterSlug,
	}

	return task, string(testData), nil
}

func parseTaskFrontmatter(raw string) (model.TaskFrontmatter, string, error) {
	const delimiter = "---"

	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, delimiter) {
		return model.TaskFrontmatter{}, "", errors.New("missing opening frontmatter delimiter")
	}

	rest := raw[len(delimiter):]
	idx := strings.Index(rest, "\n"+delimiter)
	if idx == -1 {
		return model.TaskFrontmatter{}, "", errors.New("missing closing frontmatter delimiter")
	}

	fmRaw := rest[:idx]
	content := strings.TrimSpace(rest[idx+len("\n"+delimiter):])

	var fm model.TaskFrontmatter
	if err := yaml.Unmarshal([]byte(fmRaw), &fm); err != nil {
		return model.TaskFrontmatter{}, "", err
	}

	return fm, content, nil
}

func stripTaskContent(tasks []model.Task) []model.Task {
	stripped := make([]model.Task, len(tasks))
	for i, t := range tasks {
		stripped[i] = model.Task{
			Slug:        t.Slug,
			Title:       t.Title,
			Difficulty:  t.Difficulty,
			Order:       t.Order,
			ChapterSlug: t.ChapterSlug,
		}
	}
	return stripped
}

func (s *TaskService) totalTasks() int {
	count := 0
	for _, ch := range s.chapters {
		count += len(ch.Tasks)
	}
	return count
}

func (s *TaskService) getSolvedSet(ctx context.Context, userID uuid.UUID) map[string]bool {
	solved, err := s.submissionRepo.GetSolvedTasks(ctx, userID)
	if err != nil {
		return nil
	}
	set := make(map[string]bool, len(solved))
	for _, st := range solved {
		set[st.ChapterSlug+"/"+st.TaskSlug] = true
	}
	return set
}

func (s *TaskService) enrichWithSolved(ctx context.Context, userID uuid.UUID, chapter *model.TaskChapter) int {
	solvedSet := s.getSolvedSet(ctx, userID)
	for i := range chapter.Tasks {
		v := solvedSet[chapter.Slug+"/"+chapter.Tasks[i].Slug]
		chapter.Tasks[i].Solved = &v
	}
	return len(solvedSet)
}
