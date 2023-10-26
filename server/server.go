package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"root/shared"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/vgarvardt/gue/v5"
)

const (
	queueName        = "main_queue"
	jobTypePrintName = "PrintName"
)

// a map of job types to handlers to run them
var workmap = gue.WorkMap{
	jobTypePrintName: HandlePrintNameJob,
}

type printNameArgs struct {
	Name string
}

func HandlePrintNameJob(ctx context.Context, j *gue.Job) error {
	var args printNameArgs
	if err := json.Unmarshal(j.Args, &args); err != nil {
		return err
	}
	fmt.Printf("Hello %s!\n", args.Name)
	return nil
}

func main() {
	logger := shared.L.Named("server")

	heartbeat := make(chan interface{})
	go func() {
		defer close(heartbeat)
		pulse := time.Tick(time.Second * 2)
		sendPulse := func(msg string) {
			select {
			case heartbeat <- msg:
			default:
			}
		}
		for {
			<-pulse
			sendPulse("💖")
		}
	}()
	go func() {
		for msg := range heartbeat {
			logger.Debug(msg.(string))
		}
	}()

	logger.Named("godotenv").Debug("Loading .env file")
	godotenv.Load()
	// ignore err — Docker Compose will set the environment variables

	shared.RunMigrations(logger)

	// create a queue client
	// The server is responsible for enqueueing jobs
	queue := shared.NewQueue("server")

	// start HTTP server on port 8080
	router := http.NewServeMux()

	db := shared.NewDB().DB
	// list jobs

	type GueJobRow struct {
		ID         string  `json:"id"`
		Priority   int     `json:"priority"`
		RunAt      string  `json:"run_at"`
		JobType    string  `json:"job_type"`
		Args       string  `json:"args"`
		ErrorCount int     `json:"error_count"`
		LastError  *string `json:"last_error"`
		Queue      string  `json:"queue"`
		CreatedAt  string  `json:"created_at"`
		UpdatedAt  string  `json:"updated_at"`
	}

	router.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		// job_id | priority | run_at | job_type | args | error_count | last_error | queue | created_at | updated_at
		res, err := db.Query("SELECT * FROM gue_jobs")
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		defer res.Close()

		var rows []GueJobRow

		for res.Next() {
			var row GueJobRow
			err := res.Scan(&row.ID, &row.Priority, &row.RunAt, &row.JobType, &row.Args, &row.ErrorCount, &row.LastError, &row.Queue, &row.CreatedAt, &row.UpdatedAt)
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			rows = append(rows, row)
		}

		response, _ := json.Marshal(rows)

		w.Write(response)
	})

	router.HandleFunc("/enqueue", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		name := query.Get("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		args, _ := json.Marshal(printNameArgs{Name: name})
		job := &gue.Job{
			Type:  jobTypePrintName,
			Queue: queueName,
			// RunAt: time.Now().UTC().Add(5 * time.Second), // delay
			Args: args,
		}

		if err := queue.Enqueue(context.Background(), job); err != nil {
			logger.Error("Failed to enqueue job", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			logger.Info("Job enqueued", "job_id", job.ID)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	// start HTTP server on port 8080
	http.ListenAndServe(":8080", router)
}
