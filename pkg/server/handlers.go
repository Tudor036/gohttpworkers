package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Tudor036/gohttpworkers/pkg/storage"
	"github.com/Tudor036/gohttpworkers/pkg/workers"
	"github.com/google/uuid"
)

func (s *Server) makeHTTPHandlerCallbackFunc(callback workers.WorkerCallback) http.HandlerFunc {
	return makeHTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		wj := workers.NewWorkerJob(nil)

		ctx := workers.NewWorkerCtx(wj, callback)
		s.Pool.Enqueue(ctx)

		message := "Job enqueued successfully"
		data := map[string]uuid.UUID{
			"job_id": wj.ID,
		}

		return SendResponse(w, http.StatusCreated, NewSuccessResponse(message, data))
	})
}

func makeHTTPHandlerFunc(apiFunc ApiFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := apiFunc(w, r); err != nil {
			log.Println(err.Error())
			SendResponse(w, http.StatusBadRequest, NewErrorResponse(err))
		}
	})
}

func DefaultEnqueueHandler(wj *workers.WorkerJob, worker uuid.UUID) {
	time.Sleep(time.Second)
	log.Println("job with id", wj.ID, "is finished")
}

func fetchJob(storage *storage.Storage, jobID uuid.UUID) (*workers.WorkerJob, error) {
	key := fmt.Sprintf("job-%s", jobID)

	serializedJob, err := storage.Client.Get(context.Background(), key).Result()

	if err != nil {
		return nil, fmt.Errorf("Falied to fetch job: %s", err.Error())
	}

	var job *workers.WorkerJob
	err = json.Unmarshal([]byte(serializedJob), &job)

	if err != nil {
		return nil, fmt.Errorf("Falied to deserialize job: %s", err.Error())
	}

	return job, nil
}

type statusCheckHandlerResponseBody struct {
	JobID uuid.UUID `json:"job_id"`
}

func (s *Server) DefaultStatusCheckHandler(w http.ResponseWriter, r *http.Request) error {
	var body statusCheckHandlerResponseBody
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		return fmt.Errorf("Failed to parse the body: %s", err.Error())
	}

	job, err := fetchJob(s.options.Storage, body.JobID)

	if err != nil {
		return err
	}

	return SendResponse(w, 200, NewSuccessResponse("", map[string]workers.WorkerJob{
		"job": *job,
	}))
}

func defaultNotFoundHandler(w http.ResponseWriter, r *http.Request) error {
	return SendResponse(w, http.StatusMethodNotAllowed, NewErrorResponse(errors.New("Route not found")))
}

func defaultMethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) error {
	return SendResponse(w, http.StatusMethodNotAllowed, NewErrorResponse(errors.New("Method not allowed")))
}

func (s *Server) bootstrapHandlers() {
	var handleFunc http.HandlerFunc

	if s.EnqueueCallback != nil {
		log.Panicln("EnqueueCallback is not set")
		handleFunc = s.makeHTTPHandlerCallbackFunc(s.EnqueueCallback)
	} else {
		handleFunc = s.makeHTTPHandlerCallbackFunc(DefaultEnqueueHandler)
	}

	enqueuePath := s.options.EnqueuePath
	s.router.HandleFunc(enqueuePath, handleFunc).Methods(http.MethodPost)

	if s.StatusCheckCallback != nil {
		log.Panicln("StatusCheckCallback is not set")
		handleFunc = makeHTTPHandlerFunc(s.StatusCheckCallback)
	} else {
		handleFunc = makeHTTPHandlerFunc(s.DefaultStatusCheckHandler)
	}

	statusPath := s.options.StatusPath
	s.router.HandleFunc(statusPath, handleFunc).Methods(http.MethodGet)

	if s.NotFoundCallback != nil {
		s.router.NotFoundHandler = makeHTTPHandlerFunc(s.NotFoundCallback)
	}

	if s.MethodNotAllowedCallback != nil {
		s.router.MethodNotAllowedHandler = makeHTTPHandlerFunc(s.MethodNotAllowedCallback)
	}
}
