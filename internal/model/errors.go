package model

import (
	"encoding/json"
	"errors"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ValidationErrorResponse represents a validation error response
type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	// Marshal data to JSON
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Write response
	_, err = w.Write(js)
	return err
}

// ReadJSON reads a JSON request body
func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576) // 1MB

	// Create decoder
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body
	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	// Check for additional data
	err = dec.Decode(&struct{}{})
	if err == nil {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}
