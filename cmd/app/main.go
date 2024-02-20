package main

import (
	"log"

	"github.com/IvanMeln1k/go-todo-app/internal/handler"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
	"github.com/IvanMeln1k/go-todo-app/internal/server"
	"github.com/IvanMeln1k/go-todo-app/internal/service"
)

func main() {
	repos := repository.NewRepository();
	services := service.NewService(repos);
	handlers := handler.NewHandler(services);

	srv := new(server.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}