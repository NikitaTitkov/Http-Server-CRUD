package entities

import "sync"

// UserInfo holds information about the user.
type UserInfo struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

// User represents a user entity with associated information.
type User struct {
	ID    int64    
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Email string   `json:"email"`
	Info  UserInfo `json:"info"`
}

// SyncMap is a thread-safe map to store users.
type SyncMap struct {
	Elements map[int64]*User
	Mutex    sync.RWMutex
}
