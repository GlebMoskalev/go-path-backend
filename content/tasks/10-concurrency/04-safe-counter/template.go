package solution

import "sync"

// SafeCounter — потокобезопасный счётчик.
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

// Inc увеличивает счётчик на 1.
func (c *SafeCounter) Inc() {
	// Напишите ваш код здесь — используйте c.mu.Lock/Unlock
}

// Value возвращает текущее значение счётчика.
func (c *SafeCounter) Value() int {
	// Напишите ваш код здесь — используйте c.mu.Lock/Unlock
	return 0
}

// CountConcurrently запускает n горутин, каждая вызывает Inc(), и возвращает итог.
func CountConcurrently(n int) int {
	// Напишите ваш код здесь
	return 0
}
