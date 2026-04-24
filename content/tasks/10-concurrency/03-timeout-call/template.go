package solution

import "time"

// WithTimeout запускает op и возвращает результат, если он пришёл раньше таймаута.
func WithTimeout(op func() int, ms int) (int, bool) {
	// Напишите ваш код здесь — используйте select и time.After
	_ = time.After
	return 0, false
}
