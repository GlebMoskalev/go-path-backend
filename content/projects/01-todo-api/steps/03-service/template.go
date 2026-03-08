package service

// Реализуйте сервисный слой для работы с задачами.
//
// Импорты: "todoapi/model", "todoapi/storage"
//
// 1. Определите структуру TodoService с полем storage *storage.TodoStorage
//
// 2. Реализуйте конструктор:
//    func NewTodoService(s *storage.TodoStorage) *TodoService
//
// 3. Реализуйте методы:
//    - Create(req model.CreateTodoRequest) model.Todo
//    - GetByID(id int) (model.Todo, error)
//    - GetAll() []model.Todo
//    - Update(id int, req model.UpdateTodoRequest) (model.Todo, error)
//    - Delete(id int) error
//
// В Update реализуйте частичное обновление:
// проверяйте каждый указатель на nil перед обновлением поля.
