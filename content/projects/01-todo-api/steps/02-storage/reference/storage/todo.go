package storage

import (
	"errors"
	"sync"
	"time"

	"todoapi/model"
)

var ErrNotFound = errors.New("todo not found")

type TodoStorage struct {
	mu     sync.RWMutex
	todos  map[int]model.Todo
	nextID int
}

func NewTodoStorage() *TodoStorage {
	return &TodoStorage{
		todos:  make(map[int]model.Todo),
		nextID: 1,
	}
}

func (s *TodoStorage) Create(todo model.Todo) model.Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo.ID = s.nextID
	todo.CreatedAt = time.Now()
	s.nextID++
	s.todos[todo.ID] = todo
	return todo
}

func (s *TodoStorage) GetByID(id int) (model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return model.Todo{}, ErrNotFound
	}
	return todo, nil
}

func (s *TodoStorage) GetAll() []model.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]model.Todo, 0, len(s.todos))
	for _, t := range s.todos {
		todos = append(todos, t)
	}
	return todos
}

func (s *TodoStorage) Update(id int, todo model.Todo) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return model.Todo{}, ErrNotFound
	}
	s.todos[id] = todo
	return todo, nil
}

func (s *TodoStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)
	return nil
}
