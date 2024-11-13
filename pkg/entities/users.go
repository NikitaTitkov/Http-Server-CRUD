package entities

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

// UserCreateRequest represents the data for creating a user. Used for swagger docs.
type UserCreateRequest struct {
    Name  string   `json:"name"`
    Age   int      `json:"age"`
    Email string   `json:"email"`
    Info  UserInfo `json:"info"`
}