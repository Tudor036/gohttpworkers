package server

import (
	"encoding/json"
	"net/http"

	"github.com/Tudor036/gohttpworkers/pkg/workers"
)

type ApiFunc = func(w http.ResponseWriter, r *http.Request) error
type StatusCheckCallbackFunc = func(wj *workers.WorkerJob) error
type EnqueueCallbackFunc = workers.WorkerCallback

type responseSuccess struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type responseError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type responseData struct {
	Success *responseSuccess
	Error   *responseError
}

func NewSuccessResponse(message string, data interface{}) *responseData {
	return &responseData{
		Success: &responseSuccess{
			Status:  "success",
			Message: message,
			Data:    data,
		},
		Error: nil,
	}
}

func NewErrorResponse(err error) *responseData {
	return &responseData{
		Error: &responseError{
			Status: "error",
			Error:  err.Error(),
		},
		Success: nil,
	}
}

func SendResponse(w http.ResponseWriter, statusCode int, rd *responseData) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if rd.Success != nil {
		return json.NewEncoder(w).Encode(rd.Success)
	}

	return json.NewEncoder(w).Encode(rd.Error)
}
