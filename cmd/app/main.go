package main

import (
	"log"
	"os"

	"github.com/IvanMeln1k/go-todo-app/internal/handler"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
	"github.com/IvanMeln1k/go-todo-app/internal/server"
	"github.com/IvanMeln1k/go-todo-app/internal/service"
	"github.com/IvanMeln1k/go-todo-app/pkg/database"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := database.NewPostgresDB(database.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		User: viper.GetString("db.user"),
		Password: os.Getenv("DB_PASS"),
		DBName: viper.GetString("db.name"),
		SSLMode: viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Fatalf("error connect to db: %s", err.Error())	
	}

	repos := repository.NewRepository(db)
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