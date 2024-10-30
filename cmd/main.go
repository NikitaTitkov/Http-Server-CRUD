package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/fatih/color"
	"github.com/go-chi/chi"
)

const (
	baceurl           = "http://localhost:8080"
	CreateUserPostfix = "/users"
)

type UserInfo struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

type User struct {
	Id    int64    `json:"id"`
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Email string   `json:"email"`
	Info  UserInfo `json:"info"`
}

type syncMap struct {
	elements map[int64]*User
	mutex    sync.RWMutex
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
	Users.mutex.RLock()
	defer Users.mutex.RUnlock()

}

var Users = &syncMap{elements: make(map[int64]*User)}

func main() {
	fmt.Println(color.GreenString("Hello, world!"))
	r := chi.NewRouter()
	r.Post(CreateUserPostfix, CreateUserHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
