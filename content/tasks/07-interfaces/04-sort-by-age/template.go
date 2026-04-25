package solution

// Person — человек с именем и возрастом.
type Person struct {
	Name string
	Age  int
}

// ByAge — слайс Person, сортируемый по возрасту.
type ByAge []Person

// Len возвращает количество элементов.
func (a ByAge) Len() int {
	// Напишите ваш код здесь
	return 0
}

// Less возвращает true, если элемент i должен стоять перед j.
func (a ByAge) Less(i, j int) bool {
	// Напишите ваш код здесь
	return false
}

// Swap меняет элементы i и j местами.
func (a ByAge) Swap(i, j int) {
	// Напишите ваш код здесь
}
