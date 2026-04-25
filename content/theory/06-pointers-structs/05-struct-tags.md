---
title: "Struct Tags"
description: "json/yaml/db теги, reflect.StructTag, encoding/json подробно"
order: 5
---

# Struct Tags

Теги структур — способ добавить метаданные к полям. Используются пакетами рефлексии для кастомного поведения при сериализации, валидации и работе с базами данных.

## Синтаксис тегов

```go
type T struct {
    Field Type `key:"value" key2:"value2"`
}
```

Тег — строковый литерал в обратных кавычках после типа поля. Формат: пространство-разделённые пары `ключ:"значение"`.

```go
type User struct {
    ID        int       `json:"id" db:"user_id"`
    FirstName string    `json:"first_name" db:"first_name"`
    LastName  string    `json:"last_name" db:"last_name"`
    Email     string    `json:"email" validate:"required,email"`
    Password  string    `json:"-" db:"password_hash"` // скрыт в JSON
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

---

## encoding/json — полное руководство

### Основные опции тегов json

```go
type Article struct {
    // Переименование поля:
    ID    int    `json:"id"`
    Title string `json:"title"`

    // Пропустить поле если нулевое значение:
    Views  int    `json:"views,omitempty"`   // пропустить если 0
    Author string `json:"author,omitempty"` // пропустить если ""

    // Никогда не включать в JSON:
    Internal string `json:"-"`

    // Без переименования, только опция:
    Published bool `json:",omitempty"`

    // Строковое представление числа (для JS-совместимости):
    LargeID int64 `json:"large_id,string"` // {"large_id": "1234567890"}
}
```

### Маршаллинг

```go
a := Article{
    ID:       1,
    Title:    "Введение в Go",
    Views:    0,  // omitempty — будет пропущено
    Author:   "Алиса",
    Internal: "секрет",  // json:"-" — не попадёт в JSON
}

data, err := json.Marshal(a)
// {"id":1,"title":"Введение в Go","author":"Алиса"}
// Views(0) и Internal пропущены
```

### Анмаршаллинг

```go
jsonStr := `{"id":2,"title":"Go Tips","views":150}`

var article Article
err := json.Unmarshal([]byte(jsonStr), &article)
// article.ID = 2, article.Title = "Go Tips", article.Views = 150
// article.Internal = "" (не было в JSON, нулевое значение)
```

### json.Decoder / json.Encoder для потоков

```go
import "encoding/json"
import "os"

// Чтение из файла/сети:
file, _ := os.Open("data.json")
defer file.Close()

decoder := json.NewDecoder(file)
var users []User
err := decoder.Decode(&users)

// Запись в файл/сеть:
encoder := json.NewEncoder(os.Stdout)
encoder.SetIndent("", "  ")
encoder.Encode(users)
```

`json.Decoder` и `json.Encoder` эффективнее `json.Marshal`/`json.Unmarshal` для больших данных — работают потоково без загрузки всего в память.

---

## Кастомная сериализация

Иногда стандартное поведение не подходит. Реализуй `json.Marshaler` и `json.Unmarshaler`:

```go
type Duration struct {
    time.Duration
}

// Кастомный маршаллинг: duration как строка "1h30m"
func (d Duration) MarshalJSON() ([]byte, error) {
    return json.Marshal(d.String())
}

// Кастомный анмаршаллинг:
func (d *Duration) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    duration, err := time.ParseDuration(s)
    if err != nil {
        return err
    }
    d.Duration = duration
    return nil
}

type Config struct {
    Timeout Duration `json:"timeout"`
}

c := Config{Timeout: Duration{30 * time.Second}}
data, _ := json.Marshal(c)
fmt.Println(string(data))  // {"timeout":"30s"}

var c2 Config
json.Unmarshal(data, &c2)
fmt.Println(c2.Timeout)  // 30s
```

---

## reflect.StructTag — парсинг тегов вручную

```go
import "reflect"

type Config struct {
    Host string `env:"HOST" default:"localhost" required:"true"`
    Port int    `env:"PORT" default:"8080"`
}

func parseEnvConfig(v interface{}) {
    t := reflect.TypeOf(v)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        tag := field.Tag

        envKey := tag.Get("env")
        defaultVal := tag.Get("default")
        required := tag.Get("required") == "true"

        envVal := os.Getenv(envKey)
        if envVal == "" {
            if required {
                log.Fatalf("обязательная переменная %s не установлена", envKey)
            }
            envVal = defaultVal
        }

        fmt.Printf("%s = %s\n", field.Name, envVal)
    }
}
```

---

## Популярные теги в экосистеме Go

```go
type Model struct {
    // encoding/json:
    ID   int `json:"id,omitempty"`

    // gopkg.in/yaml.v3:
    Name string `yaml:"name"`

    // database/sql + sqlx:
    Email string `db:"email"`

    // gorm:
    Age int `gorm:"column:age;not null;default:18"`

    // go-playground/validator:
    Phone string `validate:"required,e164"`  // E.164 формат телефона

    // mapstructure (для Viper конфига):
    Config string `mapstructure:"config_path"`

    // protobuf:
    Data []byte `protobuf:"bytes,1,opt,name=data,proto3"`
}
```

---

## Типичные ошибки

**Ошибка 1**: Неэкспортированное поле с json тегом — тег игнорируется.

```go
type Bad struct {
    id int `json:"id"`  // ТИХАЯ ОШИБКА: поле неэкспортировано!
}

b := Bad{id: 42}
data, _ := json.Marshal(b)
fmt.Println(string(data))  // {} — поле не сериализовано!
```

Все поля для JSON должны начинаться с заглавной буквы.

**Ошибка 2**: Опечатка в теге — Go не предупреждает.

```go
type Bad struct {
    Name string `json: "name"` // ПРОБЕЛ после двоеточия = неверный тег
}
// json.Marshal вернёт {"Name":"..."} а не {"name":"..."}
```

Используй линтеры (`staticcheck`, `golangci-lint`) для проверки тегов.

**Ошибка 3**: Забыть передать указатель в Unmarshal.

```go
var u User
json.Unmarshal(data, u)   // ОШИБКА: нужен указатель
json.Unmarshal(data, &u)  // OK
```

---

## Итог

- Теги — метаданные в обратных кавычках после типа поля
- `json:"name"` — имя поля в JSON
- `json:"name,omitempty"` — пропустить если нулевое значение
- `json:"-"` — никогда не включать
- Неэкспортированные поля (строчные) не сериализуются — всегда ошибка
- `json.Decoder`/`json.Encoder` — для потоковой обработки больших данных
- Для кастомного поведения — реализуй `MarshalJSON()`/`UnmarshalJSON()`
- Всегда передавай указатель в `json.Unmarshal`
