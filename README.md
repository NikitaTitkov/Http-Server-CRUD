# Проект на Go

Этот проект представляет собой HTTP сервер на Go с CRUD функционалом для управления пользователями. Сервер позволяет добавлять, получать, обновлять и удалять пользователей с использованием REST API.

## Установка линтера

### Для установки линтера используйте следующую команду:

```bash
make install-golangci-lint
```

### Для запуска линтера выполните:

```bash
make lint
```

### Для запуска сервера:

```bash
go run cmd/main.go
```
### Чтобы создать исполняемый файл и запустить его, выполните:

```bash
go build cmd/main.go
./main
```
## Примеры использования API

### Добавление нового пользователя

```bash
curl -i -X POST -H "Content-Type: application/json" -d '{"id":1,"name":"Alexy Laiho","age":41,"email":"alexycobhc@example.com","info":{"street":"123 Main St","city":"Anytown"}}' http://localhost:8080/newuser
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