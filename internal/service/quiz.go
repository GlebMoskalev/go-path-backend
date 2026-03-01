package service

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"path/filepath"
	"sort"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var ErrQuestionNotFound = errors.New("question not found")

type QuizService struct {
	log       *zap.Logger
	questions map[string][]model.QuizQuestion
	allByID   map[string]model.QuizQuestion
	chapters  []model.QuizChapterInfo
}

func NewQuizService(fsys fs.FS, root string, log *zap.Logger) (*QuizService, error) {
	s := &QuizService{
		log:       log,
		questions: make(map[string][]model.QuizQuestion),
		allByID:   make(map[string]model.QuizQuestion),
	}

	if err := s.load(fsys, root); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *QuizService) ListChapters() []model.QuizChapterInfo {
	return s.chapters
}

func (s *QuizService) GetQuestions(chapterSlugs []string, limit int) []model.QuizQuestion {
	var pool []model.QuizQuestion
	if len(chapterSlugs) == 0 {
		// все главы
		for _, qs := range s.questions {
			pool = append(pool, qs...)
		}
	} else {
		for _, slug := range chapterSlugs {
			if qs, ok := s.questions[slug]; ok {
				pool = append(pool, qs...)
			}
		}
	}

	rand.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	if limit > 0 && limit < len(pool) {
		pool = pool[:limit]
	}

	return pool
}

func (s *QuizService) CheckAnswer(questionID string, answer int) (model.QuizAnswerResponse, error) {
	q, ok := s.allByID[questionID]
	if !ok {
		return model.QuizAnswerResponse{}, ErrQuestionNotFound
	}

	return model.QuizAnswerResponse{
		Correct:       answer == q.Answer,
		CorrectAnswer: q.Answer,
		Explanation:   q.Explanation,
	}, nil
}

func (s *QuizService) load(fsys fs.FS, root string) error {
	dirs, err := fs.ReadDir(fsys, root)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		chapterSlug := dir.Name()
		chapterPath := filepath.Join(root, chapterSlug)

		metaData, err := fs.ReadFile(fsys, filepath.Join(chapterPath, "meta.yaml"))
		if err != nil {
			continue
		}

		var meta struct {
			Title string `yaml:"title"`
			Order int    `yaml:"order"`
		}

		if err := yaml.Unmarshal(metaData, &meta); err != nil {
			continue
		}

		quizData, err := fs.ReadFile(fsys, filepath.Join(chapterPath, "quiz.yaml"))
		if err != nil {
			continue
		}

		var raw struct {
			Questions []model.QuizQuestion `yaml:"questions"`
		}
		if err := yaml.Unmarshal(quizData, &raw); err != nil {
			s.log.Warn("skipping quiz", zap.String("chapter", chapterSlug), zap.Error(err))
			continue
		}

		for i := range raw.Questions {
			raw.Questions[i].ID = fmt.Sprintf("%s:%d", chapterSlug, i)
			raw.Questions[i].ChapterSlug = chapterSlug
			s.allByID[raw.Questions[i].ID] = raw.Questions[i]
		}

		s.questions[chapterSlug] = raw.Questions
		s.chapters = append(s.chapters, model.QuizChapterInfo{
			Slug:          chapterSlug,
			Title:         meta.Title,
			QuestionCount: len(raw.Questions),
		})
	}

	sort.Slice(s.chapters, func(i, j int) bool {
		return s.chapters[i].Slug < s.chapters[j].Slug
	})

	s.log.Info("quiz loaded",
		zap.Int("chapters", len(s.chapters)),
		zap.Int("total_questions", len(s.allByID)),
	)

	return nil
}
