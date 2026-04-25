package server

// Реализуйте роутер и запуск HTTP-сервера.
//
// Импорты:
//   import (
//       "encoding/json"
//       "net/http"
//       "taskmanager/handler"
//   )
//
// 1. func NewRouter(h *handler.TaskHandler) http.Handler
//    Зарегистрируй маршруты в http.NewServeMux():
//    GET    /health        → {"status":"ok"}
//    POST   /tasks         → h.Create
//    GET    /tasks         → h.List
//    GET    /tasks/{id}    → h.GetByID
//    PUT    /tasks/{id}    → h.Update
//    DELETE /tasks/{id}    → h.Delete
//    Верни mux как http.Handler (не *http.ServeMux)
//
// 2. func Run(addr string, handler http.Handler) error
//    Просто: return http.ListenAndServe(addr, handler)
