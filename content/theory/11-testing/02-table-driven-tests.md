---
title: "Табличные тесты"
description: "Table-driven тесты, t.Run субтесты, t.Parallel"
order: 2
---

# Табличные тесты

Табличные тесты (table-driven tests) — идиоматичный паттерн Go для проверки функции на множестве входных данных.

## Базовый паттерн

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"оба положительных", 2, 3, 5},
        {"оба отрицательных", -1, -2, -3},
        {"положительный и отрицательный", 5, -3, 2},
        {"нули", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

`t.Run` запускает подтест с именем. Каждый case — отдельный subtest.

---

## t.Run — субтесты

`t.Run(name, func)` создаёт именованный подтест:

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        a, b    float64
        want    float64
        wantErr bool
    }{
        {
            name: "нормальное деление",
            a:    10, b: 2,
            want:    5,
            wantErr: false,
        },
        {
            name:    "деление на ноль",
            a:       10, b: 0,
            want:    0,
            wantErr: true,
        },
        {
            name: "дробный результат",
            a:    1, b: 3,
            want:    0.3333333333333333,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)

            if tt.wantErr {
                if err == nil {
                    t.Fatal("ожидалась ошибка, получили nil")
                }
                return  // если ожидалась ошибка — на этом всё
            }

            if err != nil {
                t.Fatalf("неожиданная ошибка: %v", err)
            }

            if got != tt.want {
                t.Errorf("Divide(%g, %g) = %g, want %g", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

### Запуск конкретного субтеста

```bash
go test -run "TestDivide/деление_на_ноль" ./...
# пробелы в имени заменяются на _ при запуске
```

---

## Вложенные субтесты

```go
func TestHTTPHandler(t *testing.T) {
    t.Run("GET", func(t *testing.T) {
        t.Run("возвращает 200 при корректном ID", func(t *testing.T) {
            // ...
        })
        t.Run("возвращает 404 при несуществующем ID", func(t *testing.T) {
            // ...
        })
    })

    t.Run("POST", func(t *testing.T) {
        t.Run("создаёт ресурс при корректных данных", func(t *testing.T) {
            // ...
        })
        t.Run("возвращает 400 при невалидных данных", func(t *testing.T) {
            // ...
        })
    })
}
```

Запуск всех GET-тестов:

```bash
go test -run "TestHTTPHandler/GET" ./...
```

---

## t.Parallel — параллельные тесты

`t.Parallel()` позволяет тесту выполняться параллельно с другими параллельными тестами:

```go
func TestAdd(t *testing.T) {
    t.Parallel()  // этот тест может работать параллельно

    // ...
}
```

В табличных тестах — параллельные субтесты:

```go
func TestProcess(t *testing.T) {
    tests := []struct {
        name  string
        input int
        want  int
    }{
        {"малое число", 1, 2},
        {"большое число", 100, 200},
        {"ноль", 0, 0},
    }

    for _, tt := range tests {
        tt := tt  // ВАЖНО: захватить переменную цикла (до Go 1.22)
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // каждый субтест параллелен
            got := slowProcess(tt.input)
            if got != tt.want {
                t.Errorf("got %d, want %d", got, tt.want)
            }
        })
    }
}
```

> В Go 1.22+ переменные цикла автоматически копируются для каждой итерации, `tt := tt` больше не нужно. Но добавить его не вредно — это явная документация намерения.

### Управление степенью параллелизма

```bash
go test -parallel 4 ./...  # максимум 4 параллельных теста
```

---

## Тестовые фикстуры из файлов

```go
// testdata/cases.json
// [{"input": "hello world", "want": 2}, {"input": "", "want": 0}]

func TestWordCount(t *testing.T) {
    data, err := os.ReadFile("testdata/cases.json")
    if err != nil {
        t.Fatal(err)
    }

    var cases []struct {
        Input string `json:"input"`
        Want  int    `json:"want"`
    }
    if err := json.Unmarshal(data, &cases); err != nil {
        t.Fatal(err)
    }

    for _, tc := range cases {
        t.Run(tc.Input, func(t *testing.T) {
            got := WordCount(tc.Input)
            if got != tc.Want {
                t.Errorf("WordCount(%q) = %d, want %d", tc.Input, got, tc.Want)
            }
        })
    }
}
```

---

## Golden files — файлы с эталонным выводом

Паттерн для тестирования вывода, который сложно описать строкой:

```go
var update = flag.Bool("update", false, "обновить golden files")

func TestRender(t *testing.T) {
    input := "<b>Hello</b>"
    got := Render(input)

    goldenPath := filepath.Join("testdata", "render.golden")

    if *update {
        // go test -update: перезаписать эталон
        os.WriteFile(goldenPath, []byte(got), 0644)
        return
    }

    want, err := os.ReadFile(goldenPath)
    if err != nil {
        t.Fatalf("не удалось прочитать golden file: %v", err)
    }

    if got != string(want) {
        t.Errorf("вывод не совпадает с эталоном\ngot:  %q\nwant: %q", got, string(want))
    }
}
```

---

## Итог

- Табличные тесты: слайс struct с входами/ожидаемым результатом, итерация через `t.Run`
- `t.Run(name, func)` создаёт подтест — изолированный, именованный, с собственным `t.Fatal`
- `t.Parallel()` — параллельное выполнение; `tt := tt` (до Go 1.22) для захвата переменной цикла
- `-run "TestName/SubName"` — запуск конкретного субтеста
- `testdata/` — файлы с тестовыми данными; golden files для сложного вывода
