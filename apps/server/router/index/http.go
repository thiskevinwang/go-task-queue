package index

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HTTP struct {
}

func NewHTTP() HTTP {
	return HTTP{}
}

// GET /
func (h *HTTP) IndexFactory() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Hello from server"))
	}
}

// GET /health
func (h *HTTP) Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Hello from server"))
}

// Note: these module level functions will not be available
// when instantiating a HTTP struct

// GET /alt structless
func IndexFactory() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Hello from server"))
	}
}

// GET /alt/health structless
func Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Hello from server"))
}

// Ensure the HTTP struct's methods implement the httprouter.Handle interface
var (
	// factory-style: `	router.GET("/", indexHttp.Index())`
	_ func() httprouter.Handle = (*HTTP)(nil).IndexFactory
	// regular-style: `	router.GET("/health", indexHttp.Health)`
	_ httprouter.Handle = (*HTTP)(nil).Health
	// structless-factory-style`	router.GET("/", index.IndexFactory())`
	_ func() httprouter.Handle = IndexFactory
	// structless-regular-style`	router.GET("/health", index.Health)`
	_ httprouter.Handle = Health
)
