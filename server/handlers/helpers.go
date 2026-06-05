package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

)

// JSONResponse defines the standardized envelope format for all API outputs.
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// ReadJSON decodes the request body payload into the destination pointer.
// Enforces security boundaries: 1MB size cap and prevents unknown fields.
func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	err := dec.Decode(&struct{}{})

	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}

// WriteJSON wraps the payload data into the JSONResponse envelope, marshals it, and writes it.
func WriteJSON(w http.ResponseWriter, status int, data any, message string) error {
	// Construct the standardized envelope structure.
	payload := JSONResponse{
		Error:   false,
		Message: message,
		Data:    data,
	}

	// Marshal first to catch encoding bugs safely before streaming raw bytes to the pipe.
	out, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	_, err = w.Write(out)
	return err
}

// WriteError produces a standardized error envelope output.
func WriteError(w http.ResponseWriter, status int, message string) error {
	payload := JSONResponse{
		Error:   true,
		Message: message,
	}

	out, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	return err
}

