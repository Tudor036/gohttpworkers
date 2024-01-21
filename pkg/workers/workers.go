package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Tudor036/gohttpworkers/pkg/storage"
	"github.com/google/uuid"
)

type Worker struct {
	ID uuid.UUID
}

type WorkerPool struct {
	count   uint
	jobs    chan WorkerCtx
	results chan WorkerJob
	init    sync.Once
	storage *storage.Storage
}

func NewWorkerPool(count uint, store *storage.Storage) *WorkerPool {
	wp := &WorkerPool{
		count:   count,
		jobs:    make(chan WorkerCtx, 100),
		results: make(chan WorkerJob, 100),
		storage: store,
	}

	wp.init.Do(func() {
		go wp.start()
	})

	return wp
}

func saveJob(storage *storage.Storage, job *WorkerJob) error {
	key := fmt.Sprintf("job-%s", job.ID)
	serializedJob, err := json.Marshal(job)

	if err != nil {
		return fmt.Errorf("failed to marshal job to JSON: %v", err)
	}

	return storage.Client.Set(context.Background(), key, serializedJob, 0).Err()
}

func worker(id uuid.UUID, wp *WorkerPool) {
	for ctx := range wp.jobs {
		ctx.Callback(ctx.Job, id)
		ctx.Job.Finish(id, nil)
		wp.results <- *ctx.Job
	}
}

func (wp *WorkerPool) start() {
	defer func() {
		close(wp.jobs)
	}()

	for i := 0; i < int(wp.count); i++ {
		go worker(uuid.New(), wp)
	}

	for result := range wp.results {
		if err := saveJob(wp.storage, &result); err != nil {
			log.Printf("Error saving job in Redis: %v", err)
		}
	}
}

func (wp *WorkerPool) Enqueue(ctx *WorkerCtx) {
	saveJob(wp.storage, ctx.Job)
	wp.jobs <- *ctx
}
