package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"shared"
	"shared/log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type Job struct {
	ID             string     `json:"id"`
	State          string     `json:"state"` // queued | waiting | running | success | error
	ExpirationTime *time.Time `json:"expiration_time"`
	RunnerID       *string    `json:"runner_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
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

	err := shared.RunMigrations(logger)
	if err != nil {
		logger.Error("Migrations failed", "err", err)
		panic(err)
	}

	// start HTTP server on port 8080
	router := httprouter.New()

	db := shared.NewDB().DB

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Hello from server"))
	})

	// #ListJobs
	router.GET("/jobs", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		logger.Info("/jobs", "method", r.Method)

		rows, err := db.Query("SELECT * FROM jobs")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer rows.Close()

		response := []Job{}
		for rows.Next() {
			var row Job
			err = rows.Scan(&row.ID, &row.State, &row.ExpirationTime, &row.RunnerID, &row.CreatedAt, &row.UpdatedAt)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			response = append(response, row)
		}

		serialized, err := json.Marshal(response)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(serialized)
	})

	// #QueueJob
	router.POST("/enqueue", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		logger.Info("/enqueue", "method", r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	// #GetJob
	// job: queued -> waiting
	router.GET("/getjob", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		params := r.URL.Query()
		runnerID := params.Get("runner_id")
		logger.Info("/getjob", "method", r.Method, "params", params)

		var job Job
		err := db.QueryRow(`UPDATE jobs
				SET state = 'waiting', runner_id = $1
				WHERE id = (
					SELECT id FROM jobs
					WHERE state = 'queued'
					LIMIT 1
				)
				RETURNING *`, runnerID).Scan(&job.ID, &job.State, &job.ExpirationTime, &job.RunnerID, &job.CreatedAt, &job.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				// if no rows, return nil
				w.Write([]byte("null"))
				return
			} else {
				w.Write([]byte(err.Error()))
				return
			}
		}

		serialized, err := json.Marshal(job)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(serialized)

	})

	// #AckJob
	// job: waiting -> running
	router.GET("/ackjob", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	})

	// #NackJob
	// job: waiting -> queued
	router.GET("/nackjob", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	})

	// #HeartbeatJob
	// job: running -> running
	// timeout:     -> error
	router.GET("/heartbeatjob", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	})

	// #CompleteJob
	// job: running -> success
	router.GET("/completejob", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	})

	// #ErrorJob
	// job: running -> error
	router.GET("/errorjob", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	})

	// #JobLog
	router.GET("/joblog/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	})

	// start HTTP server on port 8080
	http.ListenAndServe(":8080", router)
}
