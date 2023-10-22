package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"root/shared"

	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"github.com/vgarvardt/gue/v5"
	"golang.org/x/sync/errgroup"
)

var banner = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFC0CB")).
	Padding(2).
	MaxWidth(60).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#7D56F4")).
	BorderLeft(true).BorderBottom(true).BorderTop(true).BorderRight(true)

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

	purple := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	pink := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFC0CB")).Bold(true)

	fmt.Println(fmt.Errorf(banner.Render(fmt.Sprintf(`Hello %s!`, args.Name))))
	fmt.Println(fmt.Errorf("%s %s", pink.Render(`(\__/)`), purple.Render(`||`)))
	fmt.Println(fmt.Errorf("%s %s", pink.Render(`(•ㅅ•)`), purple.Render(`||`)))
	fmt.Println(fmt.Errorf("%s %s", pink.Render(`/`), pink.Render(`　  づ`)))

	return nil
}

func main() {
	logger := shared.L.Named("runner")
	godotenv.Load()

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
			sendPulse("heartbeat 💖")
		}
	}()
	go func() {
		for msg := range heartbeat {
			logger.Debug(msg.(string))
		}
	}()

	// create a queue client
	// The runner is responsible for picking up jobs from the queue
	queue := shared.NewQueue("runner")

	// create a single worker, rather than a pool of workers. The intent is that
	// we'll have dedicated containers-per-worker, so we don't need to run multiple
	// workers in a single container.
	worker, err := gue.NewWorker(queue.Client, workmap,
		gue.WithWorkerID("worker-1"),
		gue.WithWorkerQueue(queueName),
		gue.WithWorkerSpanWorkOneNoJob(true),
		gue.WithWorkerLogger(shared.NewLogger("worker-1", shared.L)),
	)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// run worker in a goroutine
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := worker.Run(gctx)
		if err != nil {
			panic(err)
		}
		return err
	})

	// send shutdown signal to worker

	// run until Ctrl+C
	<-gctx.Done()

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
