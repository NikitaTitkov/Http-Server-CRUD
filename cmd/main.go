package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/go-chi/chi"
)

// Defines the endpoint.
const (
	baseurl             = "localhost:8080"
	CreateUserPostfix   = "/newuser"
	GetUsersDyIDPostfix = "/users/{id}"
	GetAllusersPostfix  = "/users"
	DeleteUserPostfix   = "/users/{id}"
	UpdateUserPostfix   = "/users/{id}"
)

// UserInfo holds information about the user.
type UserInfo struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

// User represents a user entity with associated information.
type User struct {
	ID    int64    `json:"id"`
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Email string   `json:"email"`
	Info  UserInfo `json:"info"`
}

type syncMap struct {
	elements map[int64]*User
	mutex    sync.RWMutex
}

// Users is a global variable that holds a map of users.
var Users = &syncMap{elements: make(map[int64]*User)}

func main() {

	r := chi.NewRouter()

	r.Post(CreateUserPostfix, CreateUserHandler)
	r.Get(GetUsersDyIDPostfix, GetUserByIDHandler)
	r.Get(GetAllusersPostfix, GetAllusersHandler)
	r.Delete(DeleteUserPostfix, DeleteUserHandler)
	r.Patch(UpdateUserPostfix, UpdateUserHandler)

	server := &http.Server{
		Addr:         baseurl,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println(color.GreenString("Server staerted!"))

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
