# Обогатитель данных

## Запуск

```go run ./cmd/main.go [-config <путь>]```

## Миграции

Взаимодействие с миграциями происходит при помощи **[migrate](https://github.com/golang-migrate/migrate)**.

## HTTP API

Основной путь `localhost:8081/api/v1`

| Метод  | Эндпоинт   | Дополнительно                      |
|--------|------------|------------------------------------|
| POST   | /user      | Создание пользователя              |
| GET    | /user      | Получение всех пользователей [^1]  |
| GET    | /user/{id} | Получение конкретного пользователя |
| PUT    | /user/{id} | Изменение конкретного пользователя |
| DELETE | /user/{id} | Удаление конкретного пользователя  |

### Создание пользователя

```curl
curl --location 'localhost:8081/api/v1/user' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Dmitriy",
    "surname": "Ushakov",
    "patronymic": "Vasilevich"
}'
```

### Получение всех пользователей

```curl
curl --location 'localhost:8081/api/v1/user'
```

### Получение пользователей с фильтром

```curl
curl --location 'localhost:8081/api/v1/user?sort_by=country&sort_order=asc'
```

### Получение конкретного пользователя

```curl
curl --location 'localhost:8081/api/v1/user/12'
```

### Изменение конкретного пользователя

```curl
curl --location --request PUT 'localhost:8081/api/v1/user/11' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Ivan",
    "surname": "Ivanov",
    "patronymic": "Ivanovich",
    "age": 21,
    "gender": "male",
    "country": "RU"
}'
```

### Удаление конкретного пользователя

```curl
curl --location --request DELETE 'localhost:8081/api/v1/user/12'
```

[^1]: дополнительные поля: `limit`, `offset`, `sort_by` и `sort_order`