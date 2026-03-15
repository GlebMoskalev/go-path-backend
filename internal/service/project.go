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
	ErrProjectNotFound     = errors.New("project not found")
	ErrProjectStepNotFound = errors.New("project step not found")
)

type ProjectService struct {
	log            *zap.Logger
	submissionRepo repository.SubmissionRepository
	projects       []model.Project
	steps          map[string]map[string]model.ProjectStep
	references     map[string]map[string]map[string]string
	tests          map[string]map[string]map[string]string
	goMods         map[string]string
	stepOrder      map[string][]string
}

func NewProjectService(fsys fs.FS, root string, log *zap.Logger, submissionRepo repository.SubmissionRepository) (*ProjectService, error) {
	s := &ProjectService{
		log:            log,
		submissionRepo: submissionRepo,
		steps:          make(map[string]map[string]model.ProjectStep),
		references:     make(map[string]map[string]map[string]string),
		tests:          make(map[string]map[string]map[string]string),
		goMods:         make(map[string]string),
		stepOrder:      make(map[string][]string),
	}
	if err := s.load(fsys, root); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, userID *uuid.UUID) []model.Project {
	result := make([]model.Project, len(s.projects))
	for i, p := range s.projects {
		result[i] = model.Project{
			Slug:        p.Slug,
			Title:       p.Title,
			Description: p.Description,
			Order:       p.Order,
			Steps:       stripStepContent(p.Steps),
		}
	}

	if userID != nil {
		solved := s.getSolvedProjectSet(ctx, *userID)
		for i := range result {
			count := 0
			for j := range result[i].Steps {
				key := result[i].Slug + "/" + result[i].Steps[j].Slug
				v := solved[key]
				result[i].Steps[j].Solved = &v
				if v {
					count++
				}
			}

			result[i].SolvedCount = count
		}
	}

	return result
}

func (s *ProjectService) GetProject(ctx context.Context, slug string, userID *uuid.UUID) (model.Project, error) {
	for _, p := range s.projects {
		if p.Slug == slug {
			result := model.Project{
				Slug:        p.Slug,
				Title:       p.Title,
				Description: p.Description,
				Order:       p.Order,
				Steps:       stripStepContent(p.Steps),
			}
			if userID != nil {
				solved := s.getSolvedProjectSet(ctx, *userID)
				count := 0
				for i := range result.Steps {
					key := slug + "/" + result.Steps[i].Slug
					v := solved[key]
					result.Steps[i].Solved = &v
					if v {
						count++
					}
				}
				result.SolvedCount = count
			}
			return result, nil
		}
	}
	return model.Project{}, ErrProjectNotFound
}

func (s *ProjectService) GetStep(ctx context.Context, projectSlug, stepSlug string, userID *uuid.UUID) (model.ProjectStep, error) {
	projectSteps, ok := s.steps[projectSlug]
	if !ok {
		return model.ProjectStep{}, ErrProjectNotFound
	}
	step, ok := projectSteps[stepSlug]
	if !ok {
		return model.ProjectStep{}, ErrProjectStepNotFound
	}
	if userID != nil {
		solved, _ := s.submissionRepo.HasSolved(ctx, *userID, projectSlug, stepSlug)
		step.Solved = &solved
	}
	return step, nil
}

func (s *ProjectService) GetFormatContext(projectSlug, stepSlug string) (*model.FormatContext, error) {
	goMod, ok := s.goMods[projectSlug]
	if !ok {
		return nil, ErrProjectNotFound
	}

	currentStep, ok := s.steps[projectSlug][stepSlug]
	if !ok {
		return nil, ErrProjectStepNotFound
	}

	files := make(map[string]string)
	for _, slug := range s.stepOrder[projectSlug] {
		if slug == stepSlug {
			continue
		}
		for path, content := range s.references[projectSlug][slug] {
			files[path] = content
		}
	}

	return &model.FormatContext{
		GoMod:    goMod,
		Files:    files,
		FilePath: currentStep.File,
	}, nil
}

func (s *ProjectService) BuildSandboxFiles(projectSlug, stepSlug, userCode string) (map[string]string, error) {
	goMod, ok := s.goMods[projectSlug]
	if !ok {
		return nil, ErrProjectNotFound
	}

	currentStep, ok := s.steps[projectSlug][stepSlug]
	if !ok {
		return nil, ErrProjectStepNotFound
	}

	files := map[string]string{
		"go.mod": goMod,
	}

	for _, slug := range s.stepOrder[projectSlug] {
		refs := s.references[projectSlug][slug]
		if slug == stepSlug {
			files[currentStep.File] = userCode
		} else {
			for path, content := range refs {
				files[path] = content
			}
		}
	}

	for path, content := range s.tests[projectSlug][stepSlug] {
		files[path] = content
	}

	return files, nil
}

func (s *ProjectService) GetStats(ctx context.Context, userID uuid.UUID) model.ProjectsStats {
	solvedSet := s.getSolvedProjectSet(ctx, userID)

	stats := model.ProjectsStats{}

	for _, p := range s.projects {
		total := len(p.Steps)
		solved := 0
		for _, step := range p.Steps {
			if solvedSet[p.Slug+"/"+step.Slug] {
				solved++
			}
		}

		stats.TotalSteps += total
		stats.SolvedSteps += solved
		stats.Projects = append(stats.Projects, model.ProjectStatsItem{
			Slug:   p.Slug,
			Title:  p.Title,
			Total:  total,
			Solved: solved,
		})
	}

	return stats
}

func (s *ProjectService) load(fsys fs.FS, root string) error {
	projectDirs, err := fs.ReadDir(fsys, root)
	if err != nil {
		return err
	}

	for _, dir := range projectDirs {
		if !dir.IsDir() {
			continue
		}
		project, err := s.loadProject(fsys, root, dir.Name())
		if err != nil {
			s.log.Warn("skipping project", zap.String("dir", dir.Name()), zap.Error(err))
			continue
		}
		s.projects = append(s.projects, project)
	}

	sort.Slice(s.projects, func(i, j int) bool {
		return s.projects[i].Order < s.projects[j].Order
	})

	s.log.Info("projects loaded",
		zap.Int("count", len(s.projects)),
		zap.Int("total_steps", s.totalSteps()),
	)

	return nil
}

func (s *ProjectService) loadProject(fsys fs.FS, root, dirName string) (model.Project, error) {
	projectPath := filepath.Join(root, dirName)

	metaData, err := fs.ReadFile(fsys, filepath.Join(projectPath, "meta.yaml"))
	if err != nil {
		return model.Project{}, err
	}
	var meta model.ProjectMeta
	if err := yaml.Unmarshal(metaData, &meta); err != nil {
		return model.Project{}, err
	}

	goModData, err := fs.ReadFile(fsys, filepath.Join(projectPath, "go.mod.tmpl"))
	if err != nil {
		return model.Project{}, err
	}
	s.goMods[dirName] = string(goModData)

	stepsPath := filepath.Join(projectPath, "steps")
	stepDirs, err := fs.ReadDir(fsys, stepsPath)
	if err != nil {
		return model.Project{}, err
	}

	s.steps[dirName] = make(map[string]model.ProjectStep)
	s.references[dirName] = make(map[string]map[string]string)
	s.tests[dirName] = make(map[string]map[string]string)

	var steps []model.ProjectStep

	for _, sd := range stepDirs {
		if !sd.IsDir() {
			continue
		}
		step, refs, testFiles, err := s.loadStep(fsys, stepsPath, sd.Name(), dirName)
		if err != nil {
			s.log.Warn("skipping step", zap.String("dir", sd.Name()), zap.Error(err))
			continue
		}
		steps = append(steps, step)
		s.steps[dirName][step.Slug] = step
		s.references[dirName][step.Slug] = refs
		s.tests[dirName][step.Slug] = testFiles
	}

	sort.Slice(steps, func(i, j int) bool {
		return steps[i].Order < steps[j].Order
	})

	orderedSlugs := make([]string, len(steps))
	for i, step := range steps {
		orderedSlugs[i] = step.Slug
	}
	s.stepOrder[dirName] = orderedSlugs

	return model.Project{
		Slug:        dirName,
		Title:       meta.Title,
		Description: meta.Description,
		Order:       meta.Order,
		Steps:       steps,
	}, nil
}

func (s *ProjectService) loadStep(fsys fs.FS, stepsPath, dirName, projectSlug string) (model.ProjectStep, map[string]string, map[string]string, error) {
	stepPath := filepath.Join(stepsPath, dirName)

	taskData, err := fs.ReadFile(fsys, filepath.Join(stepPath, "task.md"))
	if err != nil {
		return model.ProjectStep{}, nil, nil, err
	}
	fm, description, err := parseStepFrontmatter(string(taskData))
	if err != nil {
		return model.ProjectStep{}, nil, nil, err
	}

	templateData, err := fs.ReadFile(fsys, filepath.Join(stepPath, "template.go"))
	if err != nil {
		return model.ProjectStep{}, nil, nil, err
	}

	refs, err := collectFiles(fsys, filepath.Join(stepPath, "reference"))
	if err != nil {
		return model.ProjectStep{}, nil, nil, err
	}

	testFiles, err := collectFiles(fsys, filepath.Join(stepPath, "tests"))
	if err != nil {
		return model.ProjectStep{}, nil, nil, err
	}

	completionsData, err := fs.ReadFile(fsys, filepath.Join(stepPath, "completions.yaml"))
	if err != nil {
		return model.ProjectStep{}, nil, nil, err
	}
	var completionsWrapper struct {
		Completions []model.Completion `yaml:"completions"`
	}
	err = yaml.Unmarshal(completionsData, &completionsWrapper)
	if err != nil {
		return model.ProjectStep{}, nil, nil, err
	}
	completions := completionsWrapper.Completions

	step := model.ProjectStep{
		Slug:        dirName,
		Title:       fm.Title,
		Description: description,
		Template:    string(templateData),
		Difficulty:  fm.Difficulty,
		Hints:       fm.Hints,
		File:        fm.File,
		Order:       fm.Order,
		ProjectSlug: projectSlug,
		Completions: completions,
	}

	if step.File == "" && len(refs) == 1 {
		for path := range refs {
			step.File = path
		}
	}

	return step, refs, testFiles, nil
}

func parseStepFrontmatter(raw string) (model.StepFrontmatter, string, error) {
	const delimiter = "---"
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, delimiter) {
		return model.StepFrontmatter{}, "", errors.New("missing opening frontmatter delimiter")
	}
	rest := raw[len(delimiter):]
	idx := strings.Index(rest, "\n"+delimiter)
	if idx == -1 {
		return model.StepFrontmatter{}, "", errors.New("missing closing frontmatter delimiter")
	}
	fmRaw := rest[:idx]
	content := strings.TrimSpace(rest[idx+len("\n"+delimiter):])

	var fm model.StepFrontmatter
	if err := yaml.Unmarshal([]byte(fmRaw), &fm); err != nil {
		return model.StepFrontmatter{}, "", err
	}
	return fm, content, nil
}

func collectFiles(fsys fs.FS, root string) (map[string]string, error) {
	files := make(map[string]string)
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(path, root+"/")
		files[relPath] = string(content)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func stripStepContent(steps []model.ProjectStep) []model.ProjectStep {
	stripped := make([]model.ProjectStep, len(steps))
	for i, s := range steps {
		stripped[i] = model.ProjectStep{
			Slug:        s.Slug,
			Title:       s.Title,
			Difficulty:  s.Difficulty,
			Order:       s.Order,
			ProjectSlug: s.ProjectSlug,
		}
	}
	return stripped
}

func (s *ProjectService) totalSteps() int {
	count := 0
	for _, p := range s.projects {
		count += len(p.Steps)
	}
	return count
}

func (s *ProjectService) getSolvedProjectSet(ctx context.Context, userID uuid.UUID) map[string]bool {
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
