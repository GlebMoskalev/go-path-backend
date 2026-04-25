---
title: "Обёртка ошибок"
description: "Оборачивает ошибки с контекстом через fmt.Errorf %w и разворачивает через errors.Unwrap"
order: 3
difficulty: medium
---

# Обёртка ошибок

Реализуйте две функции, демонстрирующие механизм wrapping ошибок (Go 1.13+).

1. `ReadConfig(path string) error` — если `path == ""`, возвращает `errors.New("пустой путь")`; иначе оборачивает эту ошибку через `fmt.Errorf("readConfig %q: %w", path, err)`. Для демонстрации — всегда возвращает ошибку, вызывая базовую ошибку через вспомогательную функцию `baseError()`.

2. `UnwrapAll(err error) []string` — разворачивает цепочку ошибок до конца и возвращает слайс строк с сообщениями каждой ошибки в цепочке (от внешней к внутренней).

В шаблоне уже определена `baseError()`.

## Пример

```go
err := ReadConfig("/etc/app.conf")
// err.Error() = `readConfig "/etc/app.conf": базовая ошибка конфига`

UnwrapAll(err)
// ["readConfig \"/etc/app.conf\": базовая ошибка конфига", "базовая ошибка конфига"]
```
