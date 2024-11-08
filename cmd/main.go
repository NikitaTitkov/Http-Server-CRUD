package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TitkovNikita/Http-Server-CRUD/pkg/handlers"
	"github.com/fatih/color"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

func main() {
	err := initConfig()
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()

	r.Post(viper.GetString("CreateUserPostfix"), handlers.CreateUserHandler)
	r.Get(viper.GetString("GetUsersDyIDPostfix"), handlers.GetUserByIDHandler)
	r.Get(viper.GetString("GetAllusersPostfix"), handlers.GetAllusersHandler)
	r.Delete(viper.GetString("DeleteUserPostfix"), handlers.DeleteUserHandler)
	r.Patch(viper.GetString("UpdateUserPostfix"), handlers.UpdateUserHandler)

	server := &http.Server{
		Addr:         viper.GetString("baseurl"),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println(color.GreenString("Server started!"))

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
