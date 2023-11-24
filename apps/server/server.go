package main

import (
	"net/http"

	"server/router"
	"shared"
	"shared/heartbeat"
	"shared/log"

	"github.com/joho/godotenv"
)

func main() {
	logger := log.New("server")

	heartbeat.StartHeartbeat(func(msg string) {
		logger.Debug("Heartbeat", "msg", msg)
	})

	logger.Named("godotenv").Info("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Error loading .env file", "err", err)
	}
	// ignore err — Docker Compose will set the environment variables

	err = shared.RunMigrations(logger)
	if err != nil {
		logger.Error("Migrations failed", "err", err)
		panic(err)
	}

	r := router.New(router.RouterDependencies{
		Logger: logger,
		DB:     shared.NewDB().DB,
	})

	// start HTTP server on port 8080
	logger.Debug("Starting HTTP server", "port", 8080)
	http.ListenAndServe(":8080", r)
}
