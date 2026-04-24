package solution

import "fmt"

// Product — товар с названием и ценой.
type Product struct {
	Name  string
	Price int
}

// String реализует интерфейс fmt.Stringer.
func (p Product) String() string {
	// Напишите ваш код здесь
	_ = fmt.Sprintf
	return ""
}
