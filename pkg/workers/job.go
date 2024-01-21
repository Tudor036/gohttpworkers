package workers

import (
	"time"

	"github.com/google/uuid"
)

type WorkerCallback = func(wj *WorkerJob, worker uuid.UUID)

type WorkerCtx struct {
	Job      *WorkerJob
	Callback WorkerCallback
	ReqBody  interface{}
}

func NewWorkerCtx(wj *WorkerJob, callback WorkerCallback) *WorkerCtx {
	return &WorkerCtx{
		Job:      wj,
		Callback: callback,
	}
}

const (
	statusPending = "pending"
	statusDone    = "done"
)

type WorkerJob struct {
	ID        uuid.UUID
	Worker    uuid.UUID
	Status    string
	Metadata  interface{}
	Payload   interface{}
	CreatedAt time.Time
}

func NewWorkerJob(metadata interface{}) *WorkerJob {
	job := &WorkerJob{
		ID:        uuid.New(),
		Worker:    uuid.Nil,
		Status:    statusPending,
		Payload:   nil,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	return job
}

func (wj *WorkerJob) Finish(worker uuid.UUID, payload interface{}) {
	wj.Worker = worker
	wj.Payload = payload
	wj.Status = statusDone
}
