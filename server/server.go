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
