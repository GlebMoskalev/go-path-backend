package service

import (
	"todoapi/model"
	"todoapi/storage"
)

type TodoService struct {
	storage *storage.TodoStorage
}

func NewTodoService(s *storage.TodoStorage) *TodoService {
	return &TodoService{storage: s}
}

func (s *TodoService) Create(req model.CreateTodoRequest) model.Todo {
	todo := model.Todo{
		Title:       req.Title,
		Description: req.Description,
	}
	return s.storage.Create(todo)
}

func (s *TodoService) GetByID(id int) (model.Todo, error) {
	return s.storage.GetByID(id)
}

func (s *TodoService) GetAll() []model.Todo {
	return s.storage.GetAll()
}

func (s *TodoService) Update(id int, req model.UpdateTodoRequest) (model.Todo, error) {
	existing, err := s.storage.GetByID(id)
	if err != nil {
		return model.Todo{}, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Done != nil {
		existing.Done = *req.Done
	}

	return s.storage.Update(id, existing)
}

func (s *TodoService) Delete(id int) error {
	return s.storage.Delete(id)
}
