package main

import (
	"context"
	"time"

	"shared/log"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	godotenv.Load()
	logger := log.New("runner")
	childLogger := logger.Named("child")
	childLogger.Debug("Hello from child")

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
			sendPulse("🏃")
		}
	}()
	go func() {
		for msg := range heartbeat {
			logger.Debug(msg.(string))
		}
	}()

	ctx := context.Background()

	// run worker in a goroutine
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Info("Starting worker")
		return nil
	})

	// send shutdown signal to worker

	// run until Ctrl+C
	<-gctx.Done()

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
