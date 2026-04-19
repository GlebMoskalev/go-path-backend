---
title: "Рефлексия"
description: "reflect.TypeOf/ValueOf, Kind vs Type, практические применения"
order: 3
---

# Рефлексия

Рефлексия позволяет программе изучать и изменять собственную структуру во время выполнения. Используй её только когда статическая типизация не справляется.

## reflect.TypeOf и reflect.ValueOf

```go
import "reflect"

x := 42
t := reflect.TypeOf(x)   // reflect.Type — описание типа
v := reflect.ValueOf(x)  // reflect.Value — контейнер со значением

fmt.Println(t)          // int
fmt.Println(t.Name())   // int
fmt.Println(t.Kind())   // int (Kind)
fmt.Println(v.Int())    // 42
```

---

## Kind vs Type

`Kind` — категория типа. `Type` — конкретный именованный тип.

```go
type Celsius float64

var temp Celsius = 36.6
t := reflect.TypeOf(temp)

fmt.Println(t.Name())   // Celsius — имя типа
fmt.Println(t.Kind())   // float64 — underlying kind
```

Все возможные `Kind`:
```
Bool, Int, Int8, Int16, Int32, Int64
Uint, Uint8, Uint16, Uint32, Uint64, Uintptr
Float32, Float64, Complex64, Complex128
Array, Chan, Func, Interface, Map, Pointer, Slice, String, Struct, UnsafePointer
```

Для ветвления всегда используй `Kind`, а не `Type`:

```go
func describe(i any) {
    v := reflect.ValueOf(i)
    switch v.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        fmt.Printf("целое: %d\n", v.Int())
    case reflect.Float32, reflect.Float64:
        fmt.Printf("вещественное: %f\n", v.Float())
    case reflect.String:
        fmt.Printf("строка: %q\n", v.String())
    case reflect.Slice:
        fmt.Printf("слайс длиной %d\n", v.Len())
    case reflect.Struct:
        fmt.Printf("структура с %d полями\n", v.NumField())
    }
}
```

---

## Работа со структурами

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age,omitempty"`
}

func inspectStruct(s any) {
    t := reflect.TypeOf(s)
    v := reflect.ValueOf(s)

    // Если передан указатель — получить элемент
    if t.Kind() == reflect.Pointer {
        t = t.Elem()
        v = v.Elem()
    }

    fmt.Printf("тип: %s, полей: %d\n", t.Name(), t.NumField())

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)   // reflect.StructField
        value := v.Field(i)   // reflect.Value

        tag := field.Tag.Get("json")
        fmt.Printf("  %s (%s) = %v [json:%s]\n",
            field.Name, field.Type, value, tag)
    }
}

p := Person{Name: "Алиса", Age: 30}
inspectStruct(p)
// тип: Person, полей: 2
//   Name (string) = Алиса [json:name]
//   Age (int) = 30 [json:age,omitempty]
```

---

## Изменение значений через рефлексию

Чтобы изменить значение, нужен указатель, и поле должно быть экспортированным:

```go
func setField(obj any, name string, value any) error {
    v := reflect.ValueOf(obj)
    if v.Kind() != reflect.Pointer || v.IsNil() {
        return fmt.Errorf("ожидается непустой указатель")
    }
    v = v.Elem()

    field := v.FieldByName(name)
    if !field.IsValid() {
        return fmt.Errorf("поле %s не найдено", name)
    }
    if !field.CanSet() {
        return fmt.Errorf("поле %s нельзя изменить (неэкспортированное)", name)
    }

    val := reflect.ValueOf(value)
    if field.Type() != val.Type() {
        return fmt.Errorf("несовместимые типы: %v != %v", field.Type(), val.Type())
    }

    field.Set(val)
    return nil
}

p := &Person{Name: "Алиса", Age: 30}
setField(p, "Age", 31)
fmt.Println(p.Age)  // 31
```

---

## Практические применения

### Маршалинг/анмаршалинг (как работает encoding/json)

```go
func marshalToMap(s any) map[string]any {
    result := make(map[string]any)

    t := reflect.TypeOf(s)
    v := reflect.ValueOf(s)

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        if !field.IsExported() {
            continue
        }

        key := field.Name
        if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
            key = strings.Split(tag, ",")[0]
        }

        result[key] = v.Field(i).Interface()
    }

    return result
}
```

### Валидация через теги

```go
type User struct {
    Name  string `validate:"required,min=2"`
    Email string `validate:"required,email"`
    Age   int    `validate:"min=18,max=120"`
}

func validate(s any) []string {
    var errors []string
    t := reflect.TypeOf(s)
    v := reflect.ValueOf(s)

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        tag := field.Tag.Get("validate")

        if strings.Contains(tag, "required") {
            if value.IsZero() {
                errors = append(errors, fmt.Sprintf("%s: обязательное поле", field.Name))
            }
        }
    }
    return errors
}
```

---

## Производительность и ограничения

Рефлексия значительно медленнее прямого кода:

```go
// Прямой вызов: ~1 ns/op
x := 42
y := x + 1

// Через рефлексию: ~100-200 ns/op
v := reflect.ValueOf(x)
y := v.Int() + 1
```

**Когда использовать рефлексию:**
- Сериализация/десериализация (encoding/json, yaml, xml)
- Валидация через теги
- Dependency injection фреймворки
- Тестовые утилиты (testify/assert использует рефлексию)

**Когда НЕ использовать:**
- В горячем пути кода (обработчики запросов, циклы)
- Когда можно решить через интерфейсы или дженерики
- Когда код становится нечитаемым

---

## Итог

- `reflect.TypeOf(v)` → `reflect.Type`: имя, kind, поля структуры, теги
- `reflect.ValueOf(v)` → `reflect.Value`: хранилище значения с методами получения/установки
- `Kind` — категория (int, struct, slice); `Type` — конкретный тип (Person, Celsius)
- Изменение через рефлексию: нужен указатель + экспортированное поле + совместимый тип
- Рефлексия медленная: кэшируй `reflect.Type` если вызываешь многократно
- Предпочитай интерфейсы и дженерики; рефлексию — только когда они не подходят
