package main

import (
	"os"

	"github.com/IvanMeln1k/go-todo-app/internal/handler"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
	"github.com/IvanMeln1k/go-todo-app/internal/server"
	"github.com/IvanMeln1k/go-todo-app/internal/service"
	"github.com/IvanMeln1k/go-todo-app/pkg/database"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := database.NewPostgresDB(database.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("DB_PASS"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error connect to db: %s", err.Error())
		return
	}

	rdb := database.NewRedisDB(database.RedisConfig{
		Host:     "127.0.0.1",
		Port:     "6380",
		DB:       0,
		Password: "redis",
	})

	repos := repository.NewRepository(db, rdb)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// func main() {
// 	rdb := database.NewRedisDB(database.RedisConfig{
// 		Host:     "127.0.0.1",
// 		Port:     "6380",
// 		DB:       0,
// 		Password: "redis",
// 	})

// 	refreshToken := "aagert5er"
// 	ctx := context.Background()
// 	userId := 1

// 	_, err := rdb.ZAdd(ctx, fmt.Sprintf("sessionsuid%d", userId),
// 		redis.Z{Score: 0, Member: refreshToken}).Result()
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}
// }
