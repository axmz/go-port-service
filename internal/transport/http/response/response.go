package response

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// TODO: implement as per RFC 7807

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		const op = "transport.http.response.JSON"
		slog.Info(fmt.Sprintf("error encoding JSON: %v", err), slog.String("op", op))
	}
}

func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, Response{Status: StatusOK, Data: data})
}

func Err(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, Response{Status: StatusError, Error: msg})
}

func InternalServerError(w http.ResponseWriter, err error) {
	const op = "transport.http.response.InternalSeverError"
	slog.Info(fmt.Sprintf("Internal server error: %v", err), slog.String("op", op))
	Err(w, http.StatusInternalServerError, "internal server error")
}

func BadRequest(w http.ResponseWriter, msg string) {
	const op = "transport.http.response.BadRequest"
	slog.Info(fmt.Sprintf("Bad Request: %v", msg), slog.String("op", op))
	Err(w, http.StatusBadRequest, msg)
}

func NotFound(w http.ResponseWriter) {
	Err(w, http.StatusNotFound, "resource not found")
}

func Text(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func Html(w http.ResponseWriter, status int, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(html))
}
