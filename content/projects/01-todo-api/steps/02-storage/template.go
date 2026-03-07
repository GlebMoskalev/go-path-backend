package storage

// Реализуйте In-Memory хранилище для задач.
//
// 1. Определите переменную ошибки:
//    var ErrNotFound = errors.New("todo not found")
//
// 2. Определите структуру TodoStorage:
//    - mu     sync.RWMutex
//    - todos  map[int]model.Todo
//    - nextID int
//
// 3. Реализуйте конструктор NewTodoStorage() *TodoStorage
//
// 4. Реализуйте методы:
//    - Create(todo model.Todo) model.Todo
//    - GetByID(id int) (model.Todo, error)
//    - GetAll() []model.Todo
//    - Update(id int, todo model.Todo) (model.Todo, error)
//    - Delete(id int) error
//
// Не забудьте использовать мьютекс для потокобезопасности!
