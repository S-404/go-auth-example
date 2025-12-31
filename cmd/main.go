package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	database "github.com/s-404/go-auth-example/infrastructure/db"
	"github.com/s-404/go-auth-example/pkg/handler"
	"github.com/s-404/go-auth-example/pkg/repository"
	"github.com/s-404/go-auth-example/pkg/server"
	"github.com/s-404/go-auth-example/pkg/service"
	"github.com/spf13/viper"
)

// @title auth-server-example App API
// @version 1.0
// @description API Server for auth-server-example Application

// @host localhost:8082
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// init config
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := database.NewDbConnection(database.Config{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Username: viper.GetString("DB_USERNAME"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSL_MODE"),
	})

	if err != nil {
		logrus.Fatalf("database connection: %s", err.Error())
	}

	// di
	repositories := repository.NewRepository(db)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	// start http server
	srv := new(server.Server)
	go func() {
		if err := srv.Run(
			viper.GetString("PORT"),
			handlers.InitRoutes(),
		); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Printf("Server listening %s...", viper.GetString("PORT"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App shutting down...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.SetConfigFile("configs/config.yaml")

	return viper.ReadInConfig()
}
