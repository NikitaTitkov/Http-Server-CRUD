package main

import (
	"fmt"
	"net/http"
	_ "os"
	"time"

	"github.com/TitkovNikita/Http-Server-CRUD/pkg/databace"
	"github.com/TitkovNikita/Http-Server-CRUD/pkg/handlers"
	"github.com/fatih/color"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	err := initConfig()
	if err != nil {
		logrus.Fatal("failed to get config: ", err)
	}

	if err := gotenv.Load(); err != nil {
		logrus.Fatal("failed to load env variables: ", err)
	}
	DB, err := databace.NewPostgresDB(
		databace.Config{
			Host:     viper.GetString("host"),
			Port:     viper.GetString("port"),
			UserName: viper.GetString("user"),
			Password: viper.GetString("password"),
			//Password: os.Getenv("DB_PASSWORD"),
			DBname:  viper.GetString("dbname"),
			SslMode: viper.GetString("sslmode"),
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
