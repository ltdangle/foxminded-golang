package rest

import (
	"encoding/json"
	"net/http"
)

// JSON Response object.
type Response struct {
	Payload interface{} `json:"payload"`
}

func NewResponse() *Response {
	return &Response{}
}

// responder interface.
type ResponderInterface interface {
	Error(w http.ResponseWriter, message string, statusCode int)
	Success(w http.ResponseWriter, payload interface{})
}

// responder.
type responder struct{}

func NewResponder() *responder {
	return &responder{}
}

func (rspndr *responder) Error(w http.ResponseWriter, message string, statusCode int) {
	// Build json response.
	r := NewResponse()
	r.Payload = message

	// Send response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(r)
}

func (rspndr *responder) Success(w http.ResponseWriter, payload interface{}) {
	// Build json response.
	r := NewResponse()
	r.Payload = payload

	// Send response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(r)
}
