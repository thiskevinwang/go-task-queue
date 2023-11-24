// router is our entrypoint for all HTTP requests to our server
// This is where we define all of our routes and their handlers
// or controllers.
package router

import (
	"database/sql"
	"net/http"

	"shared/log"

	"server/router/index"
	"server/router/jobs"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	*httprouter.Router
	db     *sql.DB
	logger log.Logger

	// http handlers
	jobs  jobs.HTTP
	index index.HTTP
}

// ensure Router implements http.Handler
var _ http.Handler = &Router{}

type RouterDependencies struct {
	DB     *sql.DB
	Logger log.Logger
}

// New returns a new Router with dependencies, but no routes
func New(deps RouterDependencies) Router {
	return Router{
		Router: httprouter.New(),
		db:     deps.DB,
		logger: deps.Logger,

		// http handlers
		index: index.NewHTTP(),
		jobs:  jobs.NewHTTP(deps.DB),
	}
}

func (router *Router) RegisterRoutes() {
	// #Index
	router.GET("/", router.index.IndexFactory())
	router.GET("/health", router.index.Health)

	// #Jobs
	router.GET("/jobs", router.jobs.ListJobs)
	router.POST("/jobs", router.jobs.ClaimJob)
	router.POST("/jobs/enqueue", router.jobs.EnqueueJob)
}
