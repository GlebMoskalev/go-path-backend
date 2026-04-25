package solution

import "math"

// Shape — фигура, умеющая вычислять свою площадь.
type Shape interface {
	Area() float64
}

// Circle — круг.
type Circle struct {
	Radius float64
}

// Area возвращает площадь круга.
func (c Circle) Area() float64 {
	// Напишите ваш код здесь
	_ = math.Pi
	return 0
}

// Rectangle — прямоугольник.
type Rectangle struct {
	Width, Height float64
}

// Area возвращает площадь прямоугольника.
func (r Rectangle) Area() float64 {
	// Напишите ваш код здесь
	return 0
}

// TotalArea возвращает суммарную площадь всех фигур.
func TotalArea(shapes []Shape) float64 {
	// Напишите ваш код здесь
	return 0
}
