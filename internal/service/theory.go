package service

import (
	"errors"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var (
	ErrChapterNotFound = errors.New("chapter not found")
	ErrLessonNotFound  = errors.New("lesson not found")
)

type TheoryService struct {
	log      *zap.Logger
	chapters []model.Chapter
	lessons  map[string]map[string]model.Lesson
}

func NewTheoryService(fsys fs.FS, root string, log *zap.Logger) (*TheoryService, error) {
	s := TheoryService{
		log:     log,
		lessons: make(map[string]map[string]model.Lesson),
	}

	if err := s.load(fsys, root); err != nil {
		return nil, err
	}

	return &s, nil
}

// ListChapters возвращает все главы со списком уроков, но БЕЗ содержимого уроков
func (s *TheoryService) ListChapters() []model.Chapter {
	result := make([]model.Chapter, len(s.chapters))
	for i, ch := range s.chapters {
		result[i] = model.Chapter{
			Slug:        ch.Slug,
			Title:       ch.Title,
			Description: ch.Description,
			Order:       ch.Order,
			Lessons:     stripContent(ch.Lessons),
		}
	}

	return result
}

// GetChapter — возвращает одну главу по её slug (например "01-basics").
// Уроки включены, но без содержимого markdown.
// Если глава не найдена — возвращает ErrChapterNotFound.
func (s *TheoryService) GetChapter(slug string) (model.Chapter, error) {
	for _, ch := range s.chapters {
		if ch.Slug == slug {
			return model.Chapter{
				Slug:        ch.Slug,
				Title:       ch.Title,
				Description: ch.Description,
				Order:       ch.Order,
				Lessons:     stripContent(ch.Lessons),
			}, nil
		}
	}

	return model.Chapter{}, ErrChapterNotFound
}

// GetLesson — возвращает один урок С содержимым markdown.
// chapterSlug — slug главы, lessonSlug — slug урока.
func (s *TheoryService) GetLesson(chapterSlug, lessonSlug string) (model.Lesson, error) {
	chapterLessons, ok := s.lessons[chapterSlug]
	if !ok {
		return model.Lesson{}, ErrChapterNotFound
	}

	lesson, ok := chapterLessons[lessonSlug]
	if !ok {
		return model.Lesson{}, ErrLessonNotFound
	}

	return lesson, nil
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
