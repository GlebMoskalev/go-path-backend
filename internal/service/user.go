package service

import (
	"context"
	"errors"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"github.com/GlebMoskalev/go-path-backend/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	log      *zap.Logger
	userRepo repository.UserRepository
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func NewUserService(log *zap.Logger, userRepo repository.UserRepository) *UserService {
	return &UserService{log: log, userRepo: userRepo}
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get user by ID", zap.Error(err))
		if errors.Is(err, repository.UserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	s.log.Debug("successfully get user by id", zap.String("id", id.String()))
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user by email", zap.Error(err))
		return nil, err
	}

	s.log.Debug("successfully get user by email", zap.String("id", email))
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, name, picture string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user by ID", zap.Error(err))
		return err
	}

	user.Name = name
	user.Picture = picture

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.log.Error("failed to update user", zap.Error(err))
	}
	return err
}

func (s *UserService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	err := s.userRepo.Delete(ctx, userID)
	if err != nil {
		s.log.Error("failed to delete user", zap.Error(err))
	}
	return err
}
