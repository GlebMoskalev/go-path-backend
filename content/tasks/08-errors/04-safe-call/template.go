package solution

// SafeCall выполняет f и перехватывает панику, возвращая её как error.
func SafeCall(f func()) (err error) {
	// Напишите ваш код здесь — используйте defer и recover
	f()
	return
}
