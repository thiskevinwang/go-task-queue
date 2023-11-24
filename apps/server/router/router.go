// router is our entrypoint for all HTTP requests to our server
// This is where we define all of our routes and their handlers
// or controllers.
package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"shared/log"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	*httprouter.Router
	db     *sql.DB
	logger log.Logger
}

// ensure Router implements http.Handler
var _ http.Handler = &Router{}

type RouterDependencies struct {
	DB     *sql.DB
	Logger log.Logger
}

func New(deps RouterDependencies) *Router {
	return &Router{
		Router: httprouter.New(),
		db:     deps.DB,
		logger: deps.Logger,
	}
}

type Job struct {
	ID             string     `json:"id"`
	State          string     `json:"state"` // queued | waiting | running | success | error
	ExpirationTime *time.Time `json:"expiration_time"`
	RunnerID       *string    `json:"runner_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func RegisterRoutes(router *Router) {
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Hello from server"))
	})

	// #ListJobs
	router.GET("/jobs", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		rows, err := router.db.Query("SELECT * FROM jobs")
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
		w.WriteHeader(http.StatusNoContent)
	})

	// #GetJob
	// job: queued -> waiting
	router.GET("/getjob", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		params := r.URL.Query()
		runnerID := params.Get("runner_id")

		var job Job
		err := router.db.QueryRow(`UPDATE jobs
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
}
