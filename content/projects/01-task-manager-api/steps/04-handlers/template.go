package handler

// Реализуйте HTTP handlers для задач.
//
// Импорты:
//   import (
//       "context"
//       "encoding/json"
//       "net/http"
//       "strconv"
//       "taskmanager/model"
//   )
//
// 1. Определите интерфейс TaskRepository:
//    type TaskRepository interface {
//        Create(ctx context.Context, task model.Task) (model.Task, error)
//        GetByID(ctx context.Context, id int64) (model.Task, error)
//        List(ctx context.Context) ([]model.Task, error)
//        Update(ctx context.Context, task model.Task) (model.Task, error)
//        Delete(ctx context.Context, id int64) error
//    }
//
// 2. Определите структуру и конструктор:
//    type TaskHandler struct { repo TaskRepository }
//    func New(repo TaskRepository) *TaskHandler
//
// 3. Реализуйте методы:
//    Create  — POST /tasks        → 201 Created   (400 если невалидно)
//    GetByID — GET  /tasks/{id}   → 200 OK        (400 bad id, 404 not found)
//    List    — GET  /tasks        → 200 OK        (всегда массив, не null)
//    Update  — PUT  /tasks/{id}   → 200 OK        (400 bad id, 404 not found)
//    Delete  — DELETE /tasks/{id} → 204 No Content (400 bad id, 404 not found)
//
// 4. Вспомогательные функции:
//    func writeJSON(w http.ResponseWriter, status int, data any)
//    func writeError(w http.ResponseWriter, status int, message string)
//    — writeError пишет {"error": message}
//    — writeJSON устанавливает Content-Type: application/json
