package solution

// generate отправляет nums в канал в отдельной горутине и закрывает его.
func generate(nums ...int) <-chan int {
	// Напишите ваш код здесь
	out := make(chan int)
	_ = out
	return nil
}

// double читает из in, удваивает значения и отправляет в выходной канал.
func double(in <-chan int) <-chan int {
	// Напишите ваш код здесь
	out := make(chan int)
	_ = out
	return nil
}

// Pipeline соединяет generate и double, возвращает результаты.
func Pipeline(nums ...int) []int {
	// Напишите ваш код здесь
	return []int{}
}
