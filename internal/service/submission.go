package service

import (
	"context"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"github.com/GlebMoskalev/go-path-backend/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubmissionService struct {
	log            *zap.Logger
	taskService    *TaskService
	projectService *ProjectService
	sandboxService *SandboxService
	submissionRepo repository.SubmissionRepository
}

func NewSubmissionService(
	log *zap.Logger,
	taskService *TaskService,
	sandboxService *SandboxService,
	submissionRepo repository.SubmissionRepository,
	projectService *ProjectService,
) *SubmissionService {
	return &SubmissionService{
		log:            log,
		taskService:    taskService,
		sandboxService: sandboxService,
		submissionRepo: submissionRepo,
		projectService: projectService,
	}
}

func (s *SubmissionService) Submit(ctx context.Context, userID uuid.UUID, chapterSlug, taskSlug, code string) (model.SubmitResult, error) {
	testFile, err := s.taskService.GetTestFile(chapterSlug, taskSlug)
	if err != nil {
		return model.SubmitResult{}, err
	}

	result := s.sandboxService.RunTask(ctx, code, testFile)

	submission := &model.Submission{
		UserID:      userID,
		ChapterSlug: chapterSlug,
		TaskSlug:    taskSlug,
		Code:        code,
		Passed:      result.Passed,
		Result:      result,
	}

	if err := s.submissionRepo.Create(ctx, submission); err != nil {
		s.log.Error("failed to save submission", zap.Error(err))
		return model.SubmitResult{}, err
	}

	return result, nil
}

func (s *SubmissionService) SubmitProject(ctx context.Context, userID uuid.UUID, projectSlug, stepSlug, code string) (model.SubmitResult, error) {
	files, err := s.projectService.BuildSandboxFiles(projectSlug, stepSlug, code)
	if err != nil {
		return model.SubmitResult{}, err
	}

	result := s.sandboxService.RunProject(ctx, files)

	submission := &model.Submission{
		UserID:      userID,
		ChapterSlug: projectSlug,
		TaskSlug:    stepSlug,
		Code:        code,
		Passed:      result.Passed,
		Result:      result,
	}

	if err := s.submissionRepo.Create(ctx, submission); err != nil {
		s.log.Error("failed to save submission", zap.Error(err))
		return model.SubmitResult{}, err
	}

	return result, nil
}
