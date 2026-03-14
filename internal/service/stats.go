package service

import (
	"context"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"github.com/google/uuid"
)

type StatsService struct {
	theory  *TheoryService
	task    *TaskService
	project *ProjectService
}

func NewStatsService(theory *TheoryService, task *TaskService, project *ProjectService) *StatsService {
	return &StatsService{
		theory:  theory,
		task:    task,
		project: project,
	}
}

func (s *StatsService) GetUserStats(ctx context.Context, userID uuid.UUID) model.UserStats {
	return model.UserStats{
		Theory:   s.theory.GetStats(ctx, userID),
		Tasks:    s.task.GetStats(ctx, userID),
		Projects: s.project.GetStats(ctx, userID),
	}
}
