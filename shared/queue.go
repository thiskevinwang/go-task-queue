package shared

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vgarvardt/gue/v5"
	"github.com/vgarvardt/gue/v5/adapter/pgxv5"
)

type Queue struct {
	Client *gue.Client
}

func NewQueue(name string) *Queue {
	// logger := NewLogger(name, L)
	pgxCfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	pgxPool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		panic(err)
	}
	// defer pgxPool.Close()

	poolAdapter := pgxv5.NewConnPool(pgxPool)

	gc, err := gue.NewClient(
		poolAdapter,
		// gue.WithClientLogger(logger),
	)
	if err != nil {
		panic(err)
	}

	return &Queue{
		Client: gc,
	}
}

func (q *Queue) Enqueue(ctx context.Context, job *gue.Job) error {
	return q.Client.Enqueue(ctx, job)
}
