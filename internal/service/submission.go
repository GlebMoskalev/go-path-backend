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
	sandboxService *SandboxService
	submissionRepo repository.SubmissionRepository
}

func NewSubmissionService(
	log *zap.Logger,
	taskService *TaskService,
	sandboxService *SandboxService,
	submissionRepo repository.SubmissionRepository,
) *SubmissionService {
	return &SubmissionService{
		log:            log,
		taskService:    taskService,
		sandboxService: sandboxService,
		submissionRepo: submissionRepo,
	}
}

func (s *SubmissionService) Submit(ctx context.Context, userID uuid.UUID, chapterSlug, taskSlug, code string) (model.SubmitResult, error) {
	testFile, err := s.taskService.GetTestFile(chapterSlug, taskSlug)
	if err != nil {
		return model.SubmitResult{}, err
	}

	result := s.sandboxService.Run(ctx, code, testFile)

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
