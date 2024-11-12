package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/TitkovNikita/Http-Server-CRUD/docs"
	"github.com/TitkovNikita/Http-Server-CRUD/pkg/databace"
	"github.com/TitkovNikita/Http-Server-CRUD/pkg/handlers"
	"github.com/fatih/color"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

// @title HTTP SERVER CRUD API
// @version 1.0
// @description API Server for CRUD Application
// @host localhost:8080
// @BasePath /
func main() {
	flag.Parse()

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatal("failed to get config: ", err)
	}

	if err := initEnv(); err != nil {
		logrus.Fatal("failed to get env: ", err)
	}

	DB, err := databace.NewPostgresDB(
		databace.Config{
			Host:     os.Getenv("HOST"),
			Port:     os.Getenv("PORT"),
			UserName: os.Getenv("USER_DB"),
			Password: os.Getenv("DB_PASSWORD"),
			DBname:   os.Getenv("DBNAME"),
			SslMode:  os.Getenv("SSLMODE"),
		},
	)

	h := handlers.Handler{DB: DB}
	if err != nil {
		logrus.Fatal("failded to initialice db: ", err)
	}

	r := chi.NewRouter()

	r.Post(viper.GetString("CreateUserPostfix"), h.CreateUserHandler)
	r.Get(viper.GetString("GetUsersDyIDPostfix"), h.GetUserByIDHandler)
	r.Get(viper.GetString("GetAllusersPostfix"), h.GetAllusersHandler)
	r.Delete(viper.GetString("DeleteUserPostfix"), h.DeleteUserHandler)
	r.Patch(viper.GetString("UpdateUserPostfix"), h.UpdateUserHandler)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         viper.GetString("baseurl"),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println(color.GreenString("Server started!"))

	err = server.ListenAndServe()
	if err != nil {
		logrus.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func initEnv() error {
	if err := gotenv.Load(configPath); err != nil {
		return err
	}
	return nil
}
