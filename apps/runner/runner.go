package main

import (
	"context"

	"shared/heartbeat"
	"shared/log"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	godotenv.Load()
	logger := log.New("runner")

	heartbeat.StartHeartbeat(func(msg string) {
		logger.Debug("Heartbeat", "msg", msg)
	})

	ctx := context.Background()

	// run worker in a goroutine
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Info("Starting worker")
		return nil
	})

	// run until Ctrl+C
	<-gctx.Done()

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
