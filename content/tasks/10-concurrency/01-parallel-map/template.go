package solution

import "sync"

// ParallelMap применяет fn к каждому элементу параллельно и возвращает результаты.
func ParallelMap(nums []int, fn func(int) int) []int {
	result := make([]int, len(nums))
	// Напишите ваш код здесь — используйте sync.WaitGroup и горутины
	_ = sync.WaitGroup{}
	return result
}
