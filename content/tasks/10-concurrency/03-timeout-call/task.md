---
title: "Вызов с таймаутом"
description: "Выполняет операцию с ограничением по времени через select и time.After"
order: 3
difficulty: medium
---

# Вызов с таймаутом

Напишите функцию `WithTimeout`, которая запускает переданную функцию `op func() int` в горутине и ждёт результат, но не дольше `ms` миллисекунд.

- Если `op` завершилась в срок → `(результат, true)`
- Если время вышло → `(0, false)`

Используйте `select` с `time.After` для реализации таймаута.

## Пример

```go
// Быстрая операция
WithTimeout(func() int { return 42 }, 1000)
// (42, true)

// Медленная операция (спит 500мс при таймауте 100мс)
WithTimeout(func() int {
    time.Sleep(500 * time.Millisecond)
    return 0
}, 100)
// (0, false)
```
