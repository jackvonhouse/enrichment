# Обогатитель данных

## Запуск

```go run ./cmd/main.go [-config путь]```

## Миграции

Взаимодействие с миграциями происходит при помощи **[migrate](https://github.com/golang-migrate/migrate)**.
Хотя имеются возможности запустить миграции непосредственно во время работы программы,
было решено сепарировать этот функционал, избавляясь от лишней зависимости.

## Фильтрация

| Параметр   | Пример                                           | Множественное использование |
|------------|--------------------------------------------------|-----------------------------|
| name       | ```?name=Д```, ```?name=Дмитрий```               | ❌                           |
| surname    | ```?surname=Уша```, ```?surname=Ушаков```        | ❌                           |
| patronymic | ```?patronymic=```, ```?patronymic=Васильевич``` | ❌                           |
| age        | ```?age=25```                                    | ❌                           |
| ageSort    | ```?ageSort=gt```, ```?ageSort=eq```             | ❌                           |
| country    | ```?country=ru```, ```?country=UA```             | ✅                           |
| gender     | ```?gender=male```, ```?gender=female```         | ✅                           |

Возможные значения для ```ageSort```:
- gt (больше чем);
- ge (больше чем или равно);
- lt (меньше чем);
- le (меньше чем или равно);
- eq (равно);
- ne (не равно).

## Postman

Все запросы экспортированы и находятся в [postman.json](postman.json).

## HTTP API

Основной путь `localhost:8081/api/v1`

| Метод  | Эндпоинт   | Дополнительно                      |
|--------|------------|------------------------------------|
| POST   | /user      | Создание пользователя              |
| GET    | /user      | Получение всех пользователей       |
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
curl --location 'http://localhost:8081/api/v1/user?age=23&ageSort=gt'
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

## GraphQL API

Основной путь `localhost:8081/api/v1/graphql`

| Метод  | Эндпоинт   | Дополнительно                      |
|--------|------------|------------------------------------|
| POST   | /user      | Создание пользователя              |
| POST   | /user      | Получение всех пользователей       |
| POST   | /user      | Получение конкретного пользователя |
| POST   | /user      | Изменение конкретного пользователя |
| POST   | /user      | Удаление конкретного пользователя  |

### Создание пользователя

```curl
curl --location 'http://localhost:8081/api/v1/graphql/user' \
--header 'Content-Type: application/json' \
--data '{"query":"mutation {\n  create(input: {\n    name: \"John\",\n    surname: \"Doe\"\n  })\n}","variables":{}}'
```

### Получение всех пользователей

```curl
curl --location 'http://localhost:8081/api/v1/graphql/user' \
--header 'Content-Type: application/json' \
--data '{"query":"query {\n  get {\n    id\n    name\n    surname\n    patronymic\n    age\n    country\n    gender\n  }\n}","variables":{}}'
```

### Получение пользователей с фильтром

```curl
curl --location 'http://localhost:8081/api/v1/graphql/user' \
--header 'Content-Type: application/json' \
--data '{"query":"query {\n  get(get: {limit: 25}, filter: {age: 31, ageSort: \"gt\"}, sort: {sortBy: \"name\", sortOrder: \"desc\"}) {\n    id\n    name\n    surname\n    patronymic\n    age\n    country\n    gender\n  }\n}","variables":{}}'
```

### Получение конкретного пользователя

```curl
curl --location 'http://localhost:8081/api/v1/graphql/user' \
--header 'Content-Type: application/json' \
--data '{"query":"query {\n  getById(id: 1) {\n    id\n    name\n    surname\n    patronymic\n    age\n    gender\n    country\n  }\n}","variables":{}}'
```

### Изменение конкретного пользователя

```curl
curl --location 'http://localhost:8081/api/v1/graphql/user' \
--header 'Content-Type: application/json' \
--data '{"query":"mutation {\n  update(input: {\n    id: 123,\n    name: \"Ivan\",\n    surname: \"Ivanov\",\n    patronymic: \"Ivanovich\",\n    age: 31,\n    country: \"CA\",\n    gender: \"MALE\"\n  })\n}","variables":{}}'
```

### Удаление конкретного пользователя

```curl
curl --location 'http://localhost:8081/api/v1/graphql/user' \
--header 'Content-Type: application/json' \
--data '{"query":"mutation {\n  delete(id: 123)\n}","variables":{}}'
```
