# Проект на Go

Добро пожаловать в проект на Go! Этот файл README содержит инструкции по установке и запуску линтера для кода.

## Установка линтера

Для установки линтера используйте следующую команду:

```bash
make install-golangci-lint
```

Для запуска:

```bash
make lint
```

Для проврки:

```bash
curl -i -X POST -H "Content-Type: application/json" -d '{"id":1,"name":"Alexy Laiho","age":41,"email":"alexycobhc@example.com","info":{"street":"123 Main St","city":"Anytown"}}' http://localhost:8080/newuser
```

```bash
curl -i -X GET http://localhost:8080/users/1
```

```bash
curl -i -X GET http://localhost:8080/users
```