package solution

import (
	"errors"
	"fmt"
)

func baseError() error {
	return errors.New("базовая ошибка конфига")
}

// ReadConfig оборачивает базовую ошибку с контекстом пути.
func ReadConfig(path string) error {
	err := baseError()
	// Напишите ваш код здесь — используйте fmt.Errorf с %w
	_ = fmt.Errorf
	_ = err
	return nil
}

// UnwrapAll разворачивает цепочку ошибок и возвращает сообщения всех ошибок.
func UnwrapAll(err error) []string {
	// Напишите ваш код здесь — используйте errors.Unwrap в цикле
	_ = errors.Unwrap
	return nil
}
