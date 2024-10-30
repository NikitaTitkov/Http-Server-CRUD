package main

import (
	"encoding/json"
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
	baceurl             = "localhost:8080"
	CreateUserPostfix   = "/users"
	GetUsersDyIDPostfix = "/users/%d"
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

// CreateUserHandler is a handler function for creating a new user.
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}

	Users.mutex.RLock()
	defer Users.mutex.RUnlock()

}

// GetUserByIDHandler is a handler function for retrieving a user by their ID.
func GetUserByIDHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Users is a global variable that holds a map of users.
var Users = &syncMap{elements: make(map[int64]*User)}

func main() {
	fmt.Println(color.GreenString("Hello, world!"))
	r := chi.NewRouter()

	// Add a handler for creating a new user.
	r.Post(CreateUserPostfix, CreateUserHandler)
	r.Get(GetUsersDyIDPostfix, GetUserByIDHandler)

	server := &http.Server{
		Addr:         baceurl,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		//IdleTimeout:  30 * time.Second, // if needed
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
