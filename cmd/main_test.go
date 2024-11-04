package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
	// Подготовка JSON тела запроса
	user := User{
		ID:    1,
		Name:  "John Doe",
		Age:   30,
		Email: "john.doe@example.com",
		Info: UserInfo{
			Street: "123 Main St",
			City:   "New York",
		},
	}
	jsonBody, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/newuser", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Создание ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызов обработчика
	CreateUserHandler(rr, req)

	// Проверка кода ответа
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Проверка ответа
	var returnedUser User
	err := json.NewDecoder(rr.Body).Decode(&returnedUser)
	assert.NoError(t, err)
	assert.Equal(t, user, returnedUser)

	// Проверка, что пользователь добавлен в структуру Users
	Users.mutex.Lock()
	defer Users.mutex.Unlock()
	addedUser, exists := Users.elements[user.ID]
	assert.True(t, exists, "Пользователь должен быть добавлен")
	assert.Equal(t, user, *addedUser)
}
