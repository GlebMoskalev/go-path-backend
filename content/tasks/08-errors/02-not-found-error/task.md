---
title: "Кастомная ошибка NotFound"
description: "Создаёт struct-ошибку NotFoundError и использует errors.Is для её распознавания"
order: 2
difficulty: medium
---

# Кастомная ошибка NotFound

Реализуйте кастомную ошибку `NotFoundError`, которая несёт дополнительные данные — имя ресурса и его ID.

Требования:
1. Структура `NotFoundError` с полями `Resource string` и `ID int`
2. Метод `Error() string` → `"<Resource> с ID <ID> не найден"`
3. Функция `FindUser(id int) error` — возвращает `nil` если `id == 1`, иначе `&NotFoundError{"пользователь", id}`
4. Функция `IsNotFound(err error) bool` — возвращает `true` если ошибка является `*NotFoundError`, используя `errors.As`

## Пример

| `id` | `FindUser(id)`                              |
|------|---------------------------------------------|
| `1`  | `nil`                                       |
| `42` | `"пользователь с ID 42 не найден"`          |
| `0`  | `"пользователь с ID 0 не найден"`           |
