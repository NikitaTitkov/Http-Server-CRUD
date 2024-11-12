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
	"github.com/sirupsen/logrus"
)

const (
	usersTable      = "users"
	usersInfoTablle = "users_info"
)

// Handler - structure to store the database connection.
type Handler struct {
	DB *sqlx.DB
}

// CreateUserHandler is a handler function for creating a new user.
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
			logrus.Printf("Failed to rollback transaction: %v", err)
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
func (h *Handler) GetAllusersHandler(w http.ResponseWriter, _ *http.Request) {

	var users []entities.User

	query := fmt.Sprintf(`
	SELECT u.id, u.name, u.age, u.email, i.street, i.city 
	FROM %s u
	JOIN %s i ON u.info_id = i.id`, usersTable, usersInfoTablle)

	rows, err := h.DB.Queryx(query)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("Error closing rows:", err)
		}
	}()

	for rows.Next() {
		user := entities.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Email, &user.Info.Street, &user.Info.City)
		if err != nil {
			http.Error(w, "Failed to scan user data", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Failed to encode user data in get users request", http.StatusInternalServerError)
		return
	}
}

// DeleteUserHandler is a handler function for deleting a user by their ID.
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, usersTable)

	result, err := h.DB.Exec(query, id)
	if err != nil {
		logrus.Printf("Failed to delete user with ID %d: %v", id, err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateUserHandler is a handler function for updating a user's details.
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	updatedUser := &entities.User{}
	err = json.NewDecoder(r.Body).Decode(updatedUser)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}()

	tx, err := h.DB.Beginx()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("Error during rollback: %v", err)
		}
	}()

	var ageQuery interface{}
	if updatedUser.Age != 0 {
		ageQuery = updatedUser.Age
	} else {
		ageQuery = nil
	}

	userQuery := `
		UPDATE users 
		SET name = COALESCE(NULLIF($1, ''), name), 
		    age = COALESCE($2, age), 
		    email = COALESCE(NULLIF($3, ''), email)
		WHERE id = $4`
	result, err := tx.Exec(userQuery, updatedUser.Name, ageQuery, updatedUser.Email, id)
	if err != nil {
		logrus.Println("Error updating user:", err, userQuery)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	if affectedRows, _ := result.RowsAffected(); affectedRows == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if updatedUser.Info.City != "" || updatedUser.Info.Street != "" {
		infoQuery := `
			UPDATE users_info 
			SET street = COALESCE(NULLIF($1, ''), street), 
			    city = COALESCE(NULLIF($2, ''), city)
			WHERE id = (SELECT info_id FROM users WHERE id = $3)`
		_, err = tx.Exec(infoQuery, updatedUser.Info.Street, updatedUser.Info.City, id)
		if err != nil {
			logrus.Println("Error updating user info:", err)
			http.Error(w, "Failed to update user info", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
