package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

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

// DeleteUserHandler is a handler function for deleting a user by their ID.
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	Users.mutex.Lock()
	defer Users.mutex.Unlock()
	_, ok := Users.elements[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	delete(Users.elements, id)
	w.WriteHeader(http.StatusNoContent)
}

// UpdateUserHandler is a handler function for updating a user's details.
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	Users.mutex.Lock()
	defer Users.mutex.Unlock()
	user, ok := Users.elements[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	updatedUser := &User{}
	err = json.NewDecoder(r.Body).Decode(updatedUser)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	UpdateUser(user, updatedUser)
	w.WriteHeader(http.StatusNoContent)
}

// UpdateUser updates only the non-zero fields of an existing user with the values from updatedUser.
func UpdateUser(existingUser *User, updatedUser *User) {
	if updatedUser.Name != "" {
		existingUser.Name = updatedUser.Name
	}
	if updatedUser.Age != 0 {
		existingUser.Age = updatedUser.Age
	}
	if updatedUser.Email != "" {
		existingUser.Email = updatedUser.Email
	}
	if updatedUser.Info.Street != "" {
		existingUser.Info.Street = updatedUser.Info.Street
	}
	if updatedUser.Info.City != "" {
		existingUser.Info.City = updatedUser.Info.City
	}
}
