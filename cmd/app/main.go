package main

import (
	"log"

	"github.com/IvanMeln1k/go-todo-app/internal/handler"
	"github.com/IvanMeln1k/go-todo-app/internal/server"
)

func main() {
	handler := new(handler.Handler)

	server := new(server.Server)
	
	if err := server.Run("8000", handler.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}