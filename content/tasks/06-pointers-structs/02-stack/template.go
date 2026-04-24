package solution

// Stack — стек целых чисел (LIFO).
type Stack struct {
	items []int
}

// Push добавляет элемент на вершину стека.
func (s *Stack) Push(val int) {
	// Напишите ваш код здесь
}

// Pop снимает элемент с вершины стека и возвращает (значение, true).
// Если стек пустой — возвращает (0, false).
func (s *Stack) Pop() (int, bool) {
	// Напишите ваш код здесь
	return 0, false
}

// IsEmpty возвращает true, если стек не содержит элементов.
func (s *Stack) IsEmpty() bool {
	// Напишите ваш код здесь
	return true
}
