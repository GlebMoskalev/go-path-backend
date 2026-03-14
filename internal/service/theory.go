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
	ErrChapterNotFound = errors.New("chapter not found")
	ErrLessonNotFound  = errors.New("lesson not found")
)

type TheoryService struct {
	log          *zap.Logger
	chapters     []model.Chapter
	lessons      map[string]map[string]model.Lesson
	progressRepo repository.TheoryProgressRepository
}

func NewTheoryService(fsys fs.FS, root string, log *zap.Logger, progressRepo repository.TheoryProgressRepository) (*TheoryService, error) {
	s := TheoryService{
		log:          log,
		lessons:      make(map[string]map[string]model.Lesson),
		progressRepo: progressRepo,
	}

	if err := s.load(fsys, root); err != nil {
		return nil, err
	}

	return &s, nil
}

// ListChapters возвращает все главы со списком уроков, но БЕЗ содержимого уроков
func (s *TheoryService) ListChapters(ctx context.Context, userID *uuid.UUID) []model.Chapter {
	result := make([]model.Chapter, len(s.chapters))

	var completed map[string]map[string]bool
	if userID != nil {
		var err error
		completed, err = s.progressRepo.GetCompletedTheories(ctx, *userID)
		if err != nil {
			s.log.Error("failed to get completed lessons", zap.Error(err))
		}
	}

	for i, ch := range s.chapters {
		result[i] = model.Chapter{
			Slug:        ch.Slug,
			Title:       ch.Title,
			Description: ch.Description,
			Order:       ch.Order,
			Lessons:     stripContent(ch.Lessons),
		}

		if userID != nil {
			total := len(ch.Lessons)
			completedCount := 0
			if completed[ch.Slug] != nil {
				completedCount = len(completed[ch.Slug])
			}

			result[i].Progress = &model.ChapterProgress{
				Total:     total,
				Completed: completedCount,
			}

			for j := range result[i].Lessons {
				isCompleted := completed[ch.Slug] != nil && completed[ch.Slug][result[i].Lessons[j].Slug]
				result[i].Lessons[j].Completed = &isCompleted
			}
		}
	}

	return result
}

// GetChapter — возвращает одну главу по её slug (например "01-basics").
// Уроки включены, но без содержимого markdown.
// Если глава не найдена — возвращает ErrChapterNotFound.
func (s *TheoryService) GetChapter(ctx context.Context, slug string, userID *uuid.UUID) (model.Chapter, error) {
	var chapter model.Chapter
	var found bool

	for _, ch := range s.chapters {
		if ch.Slug == slug {
			chapter = model.Chapter{
				Slug:        ch.Slug,
				Title:       ch.Title,
				Description: ch.Description,
				Order:       ch.Order,
				Lessons:     stripContent(ch.Lessons),
			}
			found = true
			break
		}
	}

	if !found {
		return chapter, ErrChapterNotFound
	}

	if userID != nil {
		completed, err := s.progressRepo.GetCompletedTheories(ctx, *userID)
		if err != nil {
			s.log.Error("failed to get completed lessons", zap.Error(err))
		} else {
			total := len(chapter.Lessons)
			completedCount := 0
			if completed[slug] != nil {
				completedCount = len(completed[slug])
			}
			chapter.Progress = &model.ChapterProgress{
				Total:     total,
				Completed: completedCount,
			}

			for i := range chapter.Lessons {
				isCompleted := completed[slug] != nil && completed[slug][chapter.Lessons[i].Slug]
				chapter.Lessons[i].Completed = &isCompleted
			}
		}
	}

	return chapter, nil
}

// GetLesson — возвращает один урок С содержимым markdown.
// chapterSlug — slug главы, lessonSlug — slug урока.
func (s *TheoryService) GetLesson(ctx context.Context, chapterSlug, lessonSlug string, userID *uuid.UUID) (model.Lesson, error) {
	chapterLessons, ok := s.lessons[chapterSlug]
	if !ok {
		return model.Lesson{}, ErrChapterNotFound
	}

	lesson, ok := chapterLessons[lessonSlug]
	if !ok {
		return model.Lesson{}, ErrLessonNotFound
	}

	if userID != nil {
		isCompleted, err := s.progressRepo.IsCompleted(ctx, *userID, chapterSlug, lessonSlug)
		if err != nil {
			s.log.Error("failed to check lesson completion", zap.Error(err))
		} else {
			lesson.Completed = &isCompleted
		}
	}

	return lesson, nil
}

// MarkLessonCompleted отмечает урок как прочитанный
func (s *TheoryService) MarkLessonCompleted(ctx context.Context, userID uuid.UUID, chapterSlug, lessonSlug string) error {
	chapterLessons, ok := s.lessons[chapterSlug]
	if !ok {
		return ErrChapterNotFound
	}

	if _, ok := chapterLessons[lessonSlug]; !ok {
		return ErrLessonNotFound
	}

	return s.progressRepo.MarkCompleted(ctx, userID, chapterSlug, lessonSlug)
}

func (s *TheoryService) GetStats(ctx context.Context, userID uuid.UUID) model.TheoryStats {
	completed, err := s.progressRepo.GetCompletedTheories(ctx, userID)
	if err != nil {
		s.log.Error("failed to get completed theories for stats", zap.Error(err))
		completed = make(map[string]map[string]bool)
	}

	stats := model.TheoryStats{}

	for _, ch := range s.chapters {
		total := len(ch.Lessons)
		completedCount := len(completed[ch.Slug])

		stats.TotalLessons += total
		stats.CompletedLessons += completedCount
		stats.Chapters = append(stats.Chapters, model.TheoryChapterStats{
			Slug:      ch.Slug,
			Title:     ch.Title,
			Total:     total,
			Completed: completedCount,
		})
	}

	return stats
}

func (s *TheoryService) load(fsys fs.FS, root string) error {
	chapterDirs, err := fs.ReadDir(fsys, root)
	if err != nil {
		return err
	}
	for _, dir := range chapterDirs {
		if !dir.IsDir() {
			continue
		}

		chapter, lessons, err := s.loadChapter(fsys, root, dir.Name())
		if err != nil {
			s.log.Warn("skipping chapter", zap.String("dir", dir.Name()), zap.Error(err))
			continue
		}

		chapter.Lessons = lessons
		s.chapters = append(s.chapters, chapter)
		s.lessons[chapter.Slug] = make(map[string]model.Lesson)
		for _, l := range lessons {
			s.lessons[chapter.Slug][l.Slug] = l
		}
	}

	sort.Slice(s.chapters, func(i, j int) bool {
		return s.chapters[i].Order < s.chapters[j].Order
	})

	s.log.Info("theory loaded", zap.Int("chapters", len(s.chapters)), zap.Int("total_lessons", s.totalLessons()))

	return nil
}

func (s *TheoryService) loadChapter(fsys fs.FS, root, dirName string) (model.Chapter, []model.Lesson, error) {
	chapterPath := filepath.Join(root, dirName)
	metaPath := filepath.Join(chapterPath, "meta.yaml")

	metaData, err := fs.ReadFile(fsys, metaPath)
	if err != nil {
		return model.Chapter{}, nil, err
	}

	var meta model.ChapterMeta
	if err := yaml.Unmarshal(metaData, &meta); err != nil {
		return model.Chapter{}, nil, err
	}

	chapter := model.Chapter{
		Slug:        dirName,
		Title:       meta.Title,
		Description: meta.Description,
		Order:       meta.Order,
	}

	files, err := fs.ReadDir(fsys, chapterPath)
	if err != nil {
		return model.Chapter{}, nil, err
	}

	var lessons []model.Lesson
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		lesson, err := s.loadLesson(fsys, chapterPath, f.Name(), dirName)
		if err != nil {
			s.log.Warn("skipping lesson", zap.String("file", f.Name()), zap.Error(err))
			continue
		}
		lessons = append(lessons, lesson)
	}

	sort.Slice(lessons, func(i, j int) bool {
		return lessons[i].Order < lessons[j].Order
	})

	return chapter, lessons, nil
}

func (s *TheoryService) loadLesson(fsys fs.FS, chapterPath, fileName, chapterSlug string) (model.Lesson, error) {
	filePath := filepath.Join(chapterPath, fileName)
	data, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		return model.Lesson{}, err
	}

	fm, content, err := parseFrontmatter(string(data))
	if err != nil {
		return model.Lesson{}, err
	}

	slug := strings.TrimSuffix(fileName, ".md")

	return model.Lesson{
		Slug:        slug,
		Title:       fm.Title,
		Description: fm.Description,
		Order:       fm.Order,
		ChapterSlug: chapterSlug,
		Content:     content,
	}, nil
}

func parseFrontmatter(raw string) (model.LessonFrontmatter, string, error) {
	const delimiter = "---"

	raw = strings.TrimSpace(raw)

	if !strings.HasPrefix(raw, delimiter) {
		return model.LessonFrontmatter{}, "", errors.New("missing opening frontmatter delimiter")
	}

	rest := raw[len(delimiter):]
	idx := strings.Index(rest, "\n"+delimiter)
	if idx == -1 {
		return model.LessonFrontmatter{}, "", errors.New("missing closing frontmatter delimiter")
	}

	fmRaw := rest[:idx]
	content := strings.TrimSpace(rest[idx+len("\n"+delimiter):])

	var fm model.LessonFrontmatter
	if err := yaml.Unmarshal([]byte(fmRaw), &fm); err != nil {
		return model.LessonFrontmatter{}, "", err
	}

	return fm, content, nil
}

func (s *TheoryService) totalLessons() int {
	count := 0
	for _, ch := range s.chapters {
		count += len(ch.Lessons)
	}

	return count
}

func stripContent(lessons []model.Lesson) []model.Lesson {
	stripped := make([]model.Lesson, len(lessons))
	for i, l := range lessons {
		stripped[i] = model.Lesson{
			Slug:        l.Slug,
			Title:       l.Title,
			Description: l.Description,
			Order:       l.Order,
			ChapterSlug: l.ChapterSlug,
		}
	}

	return stripped
}
