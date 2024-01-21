package server

import "github.com/Tudor036/gohttpworkers/pkg/storage"

type ServerOptions struct {
	EnqueuePath  string
	StatusPath   string
	WorkersCount uint
	Storage      *storage.Storage
}

func DefaultServerOptions(opts ...func(*ServerOptions)) *ServerOptions {
	defaultOpts := &ServerOptions{
		EnqueuePath:  "/enqueue",
		StatusPath:   "/status",
		WorkersCount: 20,
	}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	return defaultOpts
}

func WithEnqueuePath(path string) func(*ServerOptions) {
	return func(o *ServerOptions) {
		o.EnqueuePath = path
	}
}

func WithStatusPath(path string) func(*ServerOptions) {
	return func(o *ServerOptions) {
		o.StatusPath = path
	}
}

func WithWorkersCount(count uint) func(*ServerOptions) {
	return func(o *ServerOptions) {
		o.WorkersCount = count
	}
}

func WithStorage(storage *storage.Storage) func(*ServerOptions) {
	return func(o *ServerOptions) {
		o.Storage = storage
	}
}
