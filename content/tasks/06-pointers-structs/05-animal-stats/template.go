package solution

import "fmt"

// Stats содержит физические характеристики.
type Stats struct {
	Speed  int
	Weight int
}

// Summary возвращает строку с характеристиками.
func (s Stats) Summary() string {
	return fmt.Sprintf("speed=%d, weight=%d", s.Speed, s.Weight)
}

// Animal — животное со встроенными характеристиками.
// Stats встроена без имени поля — это делает методы Stats доступными напрямую через Animal.
type Animal struct {
	Name string
	Stats
}

// Describe возвращает строку вида "<Name>: <Summary()>".
// Подсказка: вызовите a.Summary() — он доступен благодаря встраиванию.
func (a Animal) Describe() string {
	// Напишите ваш код здесь
	return ""
}
