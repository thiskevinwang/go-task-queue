package main

import (
	"encoding/json"
	"net/http"
	"time"

	"root/shared"
	"root/shared/log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

type Job struct {
	ID             string    `json:"id"`
	State          string    `json:"state"` // queued | waiting | running | success | error
	ExpirationTime time.Time `json:"expiration_time"`
	RunnerID       string    `json:"runner_id"`
}

func main() {
	logger := log.New("server")

	heartbeat := make(chan interface{})
	go func() {
		defer close(heartbeat)
		pulse := time.Tick(time.Second * 10)
		sendPulse := func(msg string) {
			select {
			case heartbeat <- msg:
			default:
			}
		}
		for {
			<-pulse
			sendPulse("🤾")
		}
	}()
	go func() {
		for msg := range heartbeat {
			logger.Debug(msg.(string))
		}
	}()

	logger.Named("godotenv").Info("Loading .env file")
	godotenv.Load()
	// ignore err — Docker Compose will set the environment variables

	shared.RunMigrations(logger)

	// start HTTP server on port 8080
	router := http.NewServeMux()

	db := shared.NewDB().DB

	router.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("/jobs", "method", r.Method)

		rows, err := db.Query("SELECT * FROM jobs")
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		defer rows.Close()

		response := []Job{}
		for rows.Next() {
			var row Job
			err = rows.Scan(&row.ID, &row.State, &row.ExpirationTime, &row.RunnerID)
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			response = append(response, row)
		}

		serialized, err := json.Marshal(response)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(serialized)
	})

	// #QueueJob
	router.HandleFunc("/enqueue", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("/enqueue", "method", r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	// For Runner to request a job
	router.HandleFunc("/getjob", func(w http.ResponseWriter, r *http.Request) {

	})

	// start HTTP server on port 8080
	http.ListenAndServe(":8080", router)
}
