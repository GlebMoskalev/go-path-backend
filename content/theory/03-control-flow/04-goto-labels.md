---
title: "Метки, goto и defer в циклах"
description: "break/continue с метками, goto (когда допустим), defer в циклах — ловушка"
order: 4
---

# Метки, goto и defer в циклах

## break и continue с метками

В Go `break` и `continue` по умолчанию работают с ближайшим внутренним циклом. Метки позволяют указать конкретный внешний цикл.

### Синтаксис меток

```go
МояМетка:
for ... {
    for ... {
        break МояМетка     // выйти из внешнего цикла
        continue МояМетка  // перейти к следующей итерации внешнего цикла
    }
}
```

### break с меткой — выход из вложенного цикла

```go
package main

import "fmt"

func main() {
    // Поиск в двумерном слайсе
    matrix := [][]int{
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9},
    }
    target := 5
    found := false

outer:
    for i, row := range matrix {
        for j, val := range row {
            if val == target {
                fmt.Printf("Нашли %d в позиции [%d][%d]\n", target, i, j)
                found = true
                break outer  // выходим из обоих циклов
            }
        }
    }

    if !found {
        fmt.Println("Не найдено")
    }
}
// Нашли 5 в позиции [1][1]
```

Без метки пришлось бы использовать флаг и дополнительный break:

```go
// Без метки — менее элегантно:
found := false
for i, row := range matrix {
    if found { break }
    for j, val := range row {
        if val == target {
            fmt.Printf("Нашли %d в [%d][%d]\n", target, i, j)
            found = true
            break
        }
    }
}
```

### continue с меткой — пропустить итерацию внешнего цикла

```go
package main

import "fmt"

func main() {
    // Найти строки, где все числа положительные
    data := [][]int{
        {1, 2, 3},
        {4, -1, 6},  // есть отрицательное
        {7, 8, 9},
    }

rows:
    for i, row := range data {
        for _, v := range row {
            if v < 0 {
                fmt.Printf("Строка %d пропущена (есть отрицательное: %d)\n", i, v)
                continue rows  // перейти к следующей строке
            }
        }
        fmt.Printf("Строка %d: все положительные\n", i)
    }
}
// Строка 0: все положительные
// Строка 1 пропущена (есть отрицательное: -1)
// Строка 2: все положительные
```

### Метки и switch

Метки работают не только с циклами, но и с `switch` и `select`. Полезно для выхода из `switch` внутри цикла:

```go
loop:
for {
    switch getCommand() {
    case "quit":
        break loop  // выходим из цикла for, а не только из switch
    case "next":
        continue loop
    default:
        process()
    }
}
```

Без метки `break` вышел бы из `switch`, но не из `for`.

---

## goto — прыжок к метке

Go поддерживает `goto`. Это редко нужно, и большинство Go-кода его не использует.

```go
goto МояМетка
// код здесь пропускается
МояМетка:
fmt.Println("перепрыгнули сюда")
```

**Ограничения `goto` в Go:**
- Нельзя прыгать через объявление переменной
- Нельзя прыгать в другой блок (только в тот же или внешний)

```go
// ОШИБКА: прыжок через объявление переменной
goto end
x := 10  // goto нельзя перепрыгнуть через это
end:
fmt.Println(x)

// Это не скомпилируется:
// goto after_decl jumps over declaration of x
```

### Когда goto допустим

В практике Go `goto` используют крайне редко. Единственный обоснованный случай — упрощение сложной логики с cleanup в низкоуровневом коде:

```go
func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }

    data, err := io.ReadAll(f)
    if err != nil {
        goto cleanup
    }

    if err = process(data); err != nil {
        goto cleanup
    }

    f.Close()
    return nil

cleanup:
    f.Close()
    return err
}
```

Однако в Go это лучше решается через `defer`:

```go
func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // гарантированно выполнится при любом возврате

    data, err := io.ReadAll(f)
    if err != nil {
        return err
    }

    return process(data)
}
```

**Вывод**: если ты пишешь `goto` — подумай, нельзя ли это решить через `defer`, функцию или ранний возврат.

---

## defer в циклах — ловушка

`defer` откладывает выполнение до момента **возврата из функции**, а не до конца текущей итерации цикла. Это приводит к классической ошибке.

### Проблема: накопление ресурсов

```go
// ОПАСНЫЙ КОД:
func processFiles(paths []string) error {
    for _, path := range paths {
        f, err := os.Open(path)
        if err != nil {
            return err
        }
        defer f.Close()  // ОШИБКА: все файлы закроются только при выходе из функции!

        if err := process(f); err != nil {
            return err
        }
    }
    return nil
}
```

Если в `paths` 1000 файлов, все они будут открыты одновременно до завершения функции. Это потенциальное исчерпание файловых дескрипторов.

### Правильное решение: вынести в отдельную функцию

```go
func processOneFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // теперь закроется при выходе из processOneFile

    return process(f)
}

func processFiles(paths []string) error {
    for _, path := range paths {
        if err := processOneFile(path); err != nil {
            return err
        }
    }
    return nil
}
```

### Альтернатива: явное закрытие

```go
func processFiles(paths []string) error {
    for _, path := range paths {
        f, err := os.Open(path)
        if err != nil {
            return err
        }

        err = process(f)
        f.Close()  // явное закрытие сразу

        if err != nil {
            return err
        }
    }
    return nil
}
```

---

## Итог

**Метки:**
- `break Метка` — выйти из помеченного цикла или switch
- `continue Метка` — перейти к следующей итерации помеченного цикла
- Используй, когда нужно управлять вложенными циклами без флагов

**goto:**
- Поддерживается, но используется крайне редко
- Нельзя прыгать через объявление переменных
- Почти всегда есть лучшее решение через `defer` или функцию

**defer в циклах:**
- `defer` выполняется при выходе из **функции**, не из итерации
- `defer` в цикле накапливает отложенные вызовы — потенциальная утечка ресурсов
- Решение: вынести тело цикла в отдельную функцию
