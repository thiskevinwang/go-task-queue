package main

import (
	"fmt"
	"net/http"

	"server/router"

	"shared"
	"shared/heartbeat"
	"shared/log"

	"github.com/joho/godotenv"
)

func main() {
	logger := log.New("server")

	heartbeat.StartHeartbeat(10, func(msg string) {
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
	r.RegisterRoutes()

	// start HTTP server on port 8080
	logger.Debug("Starting HTTP server", "port", 8080)

	// apply middleware
	wrapped := &Middleware{
		handler: r,
		logger:  log.New("middleware"),
	}

	http.ListenAndServe(":8080", wrapped)
}

type Middleware struct {
	logger  log.Logger
	handler http.Handler
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.logger.Debug(fmt.Sprintf("<-- %s %s", r.Method, r.URL.Path))
	m.handler.ServeHTTP(w, r)
}
