package main

import (
	"log"

	"github.com/IvanMeln1k/go-todo-app/internal/handler"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
	"github.com/IvanMeln1k/go-todo-app/internal/server"
	"github.com/IvanMeln1k/go-todo-app/internal/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}