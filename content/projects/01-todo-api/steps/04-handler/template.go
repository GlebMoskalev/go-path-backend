package handler

// Реализуйте HTTP обработчики для Todo API.
//
// Импорты: "todoapi/model", "todoapi/service", "todoapi/storage"
//
// 1. Определите структуру TodoHandler с полем service *service.TodoService
//
// 2. Реализуйте конструктор:
//    func NewTodoHandler(svc *service.TodoService) *TodoHandler
//
// 3. Реализуйте метод Routes() http.Handler:
//    - GET    /todos      → GetAll
//    - POST   /todos      → Create
//    - GET    /todos/{id} → GetByID
//    - PUT    /todos/{id} → Update
//    - DELETE /todos/{id} → Delete
//
// 4. Реализуйте методы-обработчики:
//    - Create(w http.ResponseWriter, r *http.Request)
//    - GetAll(w http.ResponseWriter, r *http.Request)
//    - GetByID(w http.ResponseWriter, r *http.Request)
//    - Update(w http.ResponseWriter, r *http.Request)
//    - Delete(w http.ResponseWriter, r *http.Request)
//
// 5. Реализуйте вспомогательные функции:
//    - writeJSON(w http.ResponseWriter, code int, v any)
//    - writeError(w http.ResponseWriter, code int, msg string)
//
// Используйте r.PathValue("id") для получения ID из пути.
// Create должен вернуть 201, Delete — 204.
// При пустом Title в Create — вернуть 400.
// При не найденной задаче — вернуть 404.
