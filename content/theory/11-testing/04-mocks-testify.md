---
title: "Моки и testify"
description: "httptest, интерфейсы для тестирования, testify/assert и testify/mock"
order: 4
---

# Моки и testify

Тестирование компонентов, зависящих от внешних систем (HTTP, БД), требует подмены реальных зависимостей тестовыми реализациями.

## Интерфейсы как основа тестируемости

Зависимость через интерфейс легко подменить в тестах:

```go
// service.go
type UserRepository interface {
    FindByID(id int) (*User, error)
    Save(u *User) error
}

type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUser(id int) (*User, error) {
    return s.repo.FindByID(id)
}
```

```go
// service_test.go
type mockUserRepo struct {
    users map[int]*User
    err   error
}

func (m *mockUserRepo) FindByID(id int) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.users[id], nil
}

func (m *mockUserRepo) Save(u *User) error {
    return m.err
}

func TestGetUser(t *testing.T) {
    repo := &mockUserRepo{
        users: map[int]*User{
            1: {ID: 1, Name: "Алиса"},
        },
    }
    svc := &UserService{repo: repo}

    user, err := svc.GetUser(1)
    if err != nil {
        t.Fatal(err)
    }
    if user.Name != "Алиса" {
        t.Errorf("got %s, want Алиса", user.Name)
    }
}

func TestGetUser_NotFound(t *testing.T) {
    repo := &mockUserRepo{users: map[int]*User{}}
    svc := &UserService{repo: repo}

    user, err := svc.GetUser(999)
    if err != nil {
        t.Fatal(err)
    }
    if user != nil {
        t.Errorf("ожидался nil, получили %v", user)
    }
}
```

---

## net/http/httptest — тестирование HTTP

### httptest.NewRecorder

Имитирует `http.ResponseWriter`:

```go
import "net/http/httptest"

func TestHandleUser(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
    w := httptest.NewRecorder()

    handleUser(w, req)  // вызываем handler напрямую

    resp := w.Result()
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("статус %d, ожидался %d", resp.StatusCode, http.StatusOK)
    }

    body, _ := io.ReadAll(resp.Body)
    var user User
    if err := json.Unmarshal(body, &user); err != nil {
        t.Fatalf("ошибка парсинга ответа: %v", err)
    }
    if user.ID != 1 {
        t.Errorf("ID %d, ожидался 1", user.ID)
    }
}
```

### httptest.NewServer

Поднимает реальный HTTP-сервер на случайном порту:

```go
func TestHTTPClient(t *testing.T) {
    // Создаём тестовый сервер
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(User{ID: 1, Name: "Алиса"})
    }))
    defer server.Close()  // останавливаем после теста

    // Используем URL тестового сервера
    client := NewUserClient(server.URL)
    user, err := client.GetUser(1)
    if err != nil {
        t.Fatal(err)
    }
    if user.Name != "Алиса" {
        t.Errorf("got %s, want Алиса", user.Name)
    }
}
```

---

## testify — популярная библиотека тестирования

```bash
go get github.com/stretchr/testify
```

### testify/assert

Читаемые проверки с информативными сообщениями:

```go
import "github.com/stretchr/testify/assert"

func TestDivide(t *testing.T) {
    result, err := Divide(10, 2)

    assert.NoError(t, err)
    assert.Equal(t, 5.0, result)
}

func TestDivideByZero(t *testing.T) {
    _, err := Divide(10, 0)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "деление на ноль")
}
```

При провале assert выводит:
```
Error Trace:    math_test.go:12
Error:          Not equal:
                expected: 5.0
                actual  : 4.9
Test:           TestDivide
```

### assert vs require

`assert` — продолжить тест после ошибки (аналог `t.Error`)
`require` — остановить тест после ошибки (аналог `t.Fatal`)

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
    user, err := GetUser(1)

    require.NoError(t, err)    // если ошибка — стоп, дальше нет смысла
    require.NotNil(t, user)    // если nil — стоп

    assert.Equal(t, "Алиса", user.Name)   // может продолжить
    assert.Equal(t, 25, user.Age)          // может продолжить
}
```

### Часто используемые методы

```go
assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)
assert.Nil(t, obj)
assert.NotNil(t, obj)
assert.Error(t, err)
assert.NoError(t, err)
assert.True(t, condition)
assert.False(t, condition)
assert.Contains(t, "hello world", "world")
assert.Len(t, slice, 3)
assert.Empty(t, slice)
assert.NotEmpty(t, slice)
assert.ElementsMatch(t, []int{1,2,3}, []int{3,1,2})  // без учёта порядка
assert.ErrorIs(t, err, ErrNotFound)
assert.ErrorAs(t, err, &target)
```

---

## testify/mock — моки с проверкой вызовов

```go
import "github.com/stretchr/testify/mock"

// Мок реализует интерфейс UserRepository
type MockUserRepo struct {
    mock.Mock
}

func (m *MockUserRepo) FindByID(id int) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) Save(u *User) error {
    args := m.Called(u)
    return args.Error(0)
}

func TestGetUser_WithMock(t *testing.T) {
    mockRepo := new(MockUserRepo)

    // Настраиваем ожидание: при вызове FindByID(1) вернуть пользователя
    mockRepo.On("FindByID", 1).Return(&User{ID: 1, Name: "Алиса"}, nil)

    svc := &UserService{repo: mockRepo}
    user, err := svc.GetUser(1)

    assert.NoError(t, err)
    assert.Equal(t, "Алиса", user.Name)

    // Проверяем что FindByID(1) был вызван ровно один раз
    mockRepo.AssertExpectations(t)
}
```

### Расширенные ожидания

```go
// Любой аргумент:
mockRepo.On("FindByID", mock.Anything).Return(nil, ErrNotFound)

// Несколько вызовов:
mockRepo.On("Save", mock.Anything).Return(nil).Times(3)

// Возвращать разные значения при повторных вызовах:
mockRepo.On("FindByID", 1).
    Return(&User{Name: "первый"}, nil).Once().
    Return(&User{Name: "второй"}, nil).Once()
```

---

## Итог

- Проектируй зависимости через интерфейсы — это делает код тестируемым
- Ручные моки: простой struct, реализующий интерфейс; подходит для большинства случаев
- `httptest.NewRecorder` — тестировать handler без поднятия сервера
- `httptest.NewServer` — тестировать HTTP-клиент против реального сервера
- `testify/assert` — читаемые проверки с информативными сообщениями об ошибках
- `testify/require` — аналог, но останавливает тест при первой ошибке
- `testify/mock` — моки с проверкой вызовов; для сложных сценариев
