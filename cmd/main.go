package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// CreateUserHandler is a handler function for creating a new user.// CreateUserHandler is a handler function for creating a new user.
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode user data in create request", http.StatusInternalServerError)
		return
	}

	Users.mutex.Lock()
	defer Users.mutex.Unlock()

	Users.elements[user.ID] = user

	for id, u := range Users.elements {
		log.Printf("User ID: %d, Name: %s, Age: %d, Email: %s, Info: %+v\n", id, u.Name, u.Age, u.Email, u.Info)
	}
}

// GetUserByIDHandler is a handler function for retrieving a user by their ID.
func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	Users.mutex.RLock()
	defer Users.mutex.RUnlock()
	user, ok := Users.elements[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode user data in get byID request", http.StatusInternalServerError)
		return
	}
}

// GetAllusersHandler is a handler function for retrieving all users.
func GetAllusersHandler(w http.ResponseWriter, _ *http.Request) {
	Users.mutex.RLock()
	defer Users.mutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(Users.elements)
	if err != nil {
		http.Error(w, "Failed to encode user data in get users request", http.StatusInternalServerError)
		return
	}
}

// Users is a global variable that holds a map of users.
var Users = &syncMap{elements: make(map[int64]*User)}

func main() {

	r := chi.NewRouter()

	// Add a handler for creating a new user.
	r.Post(CreateUserPostfix, CreateUserHandler)
	r.Get(GetUsersDyIDPostfix, GetUserByIDHandler)
	r.Get(GetAllusersPostfix, GetAllusersHandler)

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
