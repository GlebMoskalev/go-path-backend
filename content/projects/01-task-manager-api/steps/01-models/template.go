package model

// Определите тип Status и его константы:
//   type Status string
//   const (
//       StatusPending    Status = "pending"
//       StatusInProgress Status = "in_progress"
//       StatusDone       Status = "done"
//   )

// Определите структуру Task с JSON-тегами:
//   type Task struct {
//       ID          int64     `json:"id"`
//       Title       string    `json:"title"`
//       Description string    `json:"description"`
//       Status      Status    `json:"status"`
//       CreatedAt   time.Time `json:"created_at"`
//       UpdatedAt   time.Time `json:"updated_at"`
//   }

// NewTask создаёт новую задачу со статусом StatusPending и текущим временем.
// func NewTask(title, description string) Task { ... }

// Validate возвращает errors.New("title is required") если Title пустой.
// func (t Task) Validate() error { ... }
