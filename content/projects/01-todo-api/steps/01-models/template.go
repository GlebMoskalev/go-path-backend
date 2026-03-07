package model

// Todo представляет задачу в приложении.
// Определите структуру с полями:
//   - ID          int       `json:"id"`
//   - Title       string    `json:"title"`
//   - Description string    `json:"description"`
//   - Done        bool      `json:"done"`
//   - CreatedAt   time.Time `json:"created_at"`

// CreateTodoRequest — тело запроса на создание задачи.
// Определите структуру с полями:
//   - Title       string `json:"title"`
//   - Description string `json:"description"`

// UpdateTodoRequest — тело запроса на обновление задачи.
// Используйте указатели для частичного обновления.
// Определите структуру с полями:
//   - Title       *string `json:"title"`
//   - Description *string `json:"description"`
//   - Done        *bool   `json:"done"`
