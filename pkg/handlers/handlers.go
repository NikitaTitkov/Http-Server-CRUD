package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TitkovNikita/Http-Server-CRUD/pkg/entities"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable      = "users"
	usersInfoTablle = "users_info"
)

// Handler - structure to store the database connection.
type Handler struct {
	DB *sqlx.DB
}

// Users contains a synchronized map of users.
var Users = &entities.SyncMap{Elements: make(map[int64]*entities.User)}

// CreateUserHandler is a handler function for creating a new user.// CreateUserHandler is a handler function for creating a new user.
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &entities.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	tx, err := h.DB.Beginx()
	if err != nil {
		http.Error(w, "Failed to begin transaction: ", http.StatusInternalServerError)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Failed to rollback transaction: %v", err)
		}
	}()

	queryUserInfo := fmt.Sprintf("INSERT INTO %s (street, city) values ($1, $2) RETURNING id", usersInfoTablle)
	var infoID int
	err = tx.QueryRowx(queryUserInfo, user.Info.Street, user.Info.City).Scan(&infoID)
	if err != nil {
		http.Error(w, "Failed to insert user info", http.StatusInternalServerError)
		return
	}

	queryUser := fmt.Sprintf("INSERT INTO %s (name, age, email, info_id) VALUES ($1, $2, $3, $4) RETURNING id", usersTable)
	err = tx.QueryRowx(queryUser, user.Name, user.Age, user.Email, infoID).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "failed to commit", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode user data in create request", http.StatusInternalServerError)
		return
	}

	Users.Mutex.Lock()
	defer Users.Mutex.Unlock()

	Users.Elements[user.ID] = user

	for id, u := range Users.Elements {
		log.Printf("User ID: %d, Name: %s, Age: %d, Email: %s, Info: %+v\n", id, u.Name, u.Age, u.Email, u.Info)
	}
}

// GetUserByIDHandler is a handler function for retrieving a user by their ID.
func (h *Handler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	user := &entities.User{}
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`
	SELECT u.id, u.name, u.age, u.email, i.street, i.city 
	FROM %s u
	JOIN %s i ON u.info_id = i.id
	WHERE u.id = $1`, usersTable, usersInfoTablle)

	row := h.DB.QueryRowx(query, id)
	err = row.Scan(&user.ID, &user.Name, &user.Age, &user.Email, &user.Info.Street, &user.Info.City)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Database error: %v", err)
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
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
	Users.Mutex.RLock()
	defer Users.Mutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(Users.Elements)
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
	Users.Mutex.Lock()
	defer Users.Mutex.Unlock()
	_, ok := Users.Elements[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	delete(Users.Elements, id)
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
	Users.Mutex.Lock()
	defer Users.Mutex.Unlock()
	user, ok := Users.Elements[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	updatedUser := &entities.User{}
	err = json.NewDecoder(r.Body).Decode(updatedUser)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	UpdateUser(user, updatedUser)
	w.WriteHeader(http.StatusNoContent)
}

// UpdateUser updates only the non-zero fields of an existing user with the values from updatedUser.
func UpdateUser(existingUser *entities.User, updatedUser *entities.User) {
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
