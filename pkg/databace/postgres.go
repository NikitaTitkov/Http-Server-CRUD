package databace

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Config - structure to store the connection settings for the database.
type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBname   string
	SslMode  string
}

// NewPostgresDB is a function for creating a new connection to a Postgres database.
// It takes a Config struct as input, which contains the connection settings for the database.
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.UserName, cfg.Password, cfg.DBname, cfg.SslMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
