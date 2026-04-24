package solution

import (
	"errors"
	"fmt"
)

// NotFoundError — ошибка "ресурс не найден" с дополнительными данными.
type NotFoundError struct {
	Resource string
	ID       int
}

// Error реализует интерфейс error.
func (e *NotFoundError) Error() string {
	// Напишите ваш код здесь
	_ = fmt.Sprintf
	return ""
}

// FindUser возвращает nil если id == 1, иначе NotFoundError.
func FindUser(id int) error {
	// Напишите ваш код здесь
	return nil
}

// IsNotFound возвращает true если err является *NotFoundError.
func IsNotFound(err error) bool {
	// Напишите ваш код здесь
	_ = errors.As
	return false
}
