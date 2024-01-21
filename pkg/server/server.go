package server

import (
	"net/http"

	"github.com/Tudor036/gohttpworkers/pkg/workers"
	"github.com/gorilla/mux"
)

type Server struct {
	options                  *ServerOptions
	router                   *mux.Router
	Pool                     *workers.WorkerPool
	EnqueueCallback          workers.WorkerCallback
	StatusCheckCallback      ApiFunc
	NotFoundCallback         ApiFunc
	MethodNotAllowedCallback ApiFunc
}

func NewServer(opts *ServerOptions) *Server {
	r := mux.NewRouter()

	r.NotFoundHandler = makeHTTPHandlerFunc(defaultNotFoundHandler)
	r.MethodNotAllowedHandler = makeHTTPHandlerFunc(defaultMethodNotAllowedHandler)

	return &Server{
		options: opts,
		router:  r,
		Pool:    workers.NewWorkerPool(opts.WorkersCount, opts.Storage),
	}
}

func (s *Server) AsHTTPHandler() http.Handler {
	s.bootstrapHandlers()
	return s.router
}

func (s *Server) ListenAndServe(addr string) error {
	s.bootstrapHandlers()
	return http.ListenAndServe(addr, s.router)
}
