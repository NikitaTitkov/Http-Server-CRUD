<h1 align="center">HttpServer-CRUD</h1>

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [API Examples](#api_examples)

## About <a name = "about"></a>

HTTP сервер на Go, реализующий CRUD (создание, чтение, обновление и удаление) для управления пользователями через REST API.

## Getting Started <a name = "getting_started"></a>

Инструкции для локального запуска, сборки и установки линтера.

### Установка линтера

Для установки линтера используйте следующую команду:

```bash
make install-golangci-lint
```

### Запуск линтера

```bash
make lint
```
### [Migrations](https://github.com/golang-migrate/migrate)
Накат миграций (up)

Замените !YOUR USER!, !YOUR PASSWORD!, и !YOUR PORT! на реальные значения для подключения к базе данных

Параметр sslmode=disable отключает SSL-соединение (например, если ваше приложение не требует шифрования)

```bash
migrate -path ./schema -database 'postgres://!YOUR USER!:!YOUR PASSWORD!@localhost:!YOUR PORT!/postgres?sslmode=disable' up
```
Откат миграций (down)

Это позволит вернуться к предыдущей версии схемы базы данных

```bash
migrate -path ./schema -database 'postgres://!YOUR USER!:!YOUR PASSWORD!@localhost:!YOUR PORT!/postgres?sslmode=disable' down
```

### Запуск сервера локально

```bash
go run cmd/server/main.go
```

### Сборка и запуск сервера

```bash
go build -o bin/server ./cmd/server/main.go
./bin/server
```

## API Examples <a name = "api_examples"></a>

### Добавление нового пользователя

```bash
curl -i -X POST -H "Content-Type: application/json" -d '{"name":"Alexy Laiho","age":41,"email":"alexycobhc@example.com","info":{"street":"123 Main St","city":"Anytown"}}' http://localhost:8080/newuser
```

### Получение информации о пользователе

```bash
curl -i -X GET http://localhost:8080/users/1
```

### Получение списка всех пользователей

```bash
curl -i -X GET http://localhost:8080/users
```

### Удаление пользователя

```bash
curl -i -X DELETE http://localhost:8080/users/1
```

### Обновление информации о пользователе

```bash
curl -i -X PATCH "http://localhost:8080/users/1" \
     -H "Content-Type: application/json" \
     -d '{
           "name": "Updated Name",
           "info": {
               "city": "Updated City"
           }
         }'
```