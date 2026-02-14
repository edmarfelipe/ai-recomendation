package internal

import (
	"fmt"
	"log/slog"
	"net/http"
)

// APIError represents an error that occurred while handling an API request.
type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIError: status=%d, message=%s", e.Status, e.Message)
}

// HandleWithError is a helper to wrap HTTP handlers that return an error.
// check what to do with the error and respond accordingly.
// if an unknown error occurs, respond with HTTP 500, with the error message.
func HandleWithError(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			if apiErr, ok := err.(*APIError); ok {
				http.Error(w, apiErr.Message, apiErr.Status)
				return
			}
			slog.Error("handler error", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
