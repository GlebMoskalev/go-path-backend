package solution

// Shape — интерфейс фигуры с методом Area.
type Shape interface {
	Area() float64
}

// Rectangle — прямоугольник.
type Rectangle struct {
	Width, Height float64
}

// Circle — круг.
type Circle struct {
	Radius float64
}

// Реализуйте метод Area() для Rectangle и Circle.

// TotalArea возвращает сумму площадей всех фигур.
func TotalArea(shapes []Shape) float64 {
	// Напишите ваш код здесь
	return 0
}
