package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// TODO: implement as per RFC 7807

func Json(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding JSON: %v", err)
	}
}

func Ok(w http.ResponseWriter, data interface{}) {
	Json(w, http.StatusOK, data)
}

func Err(w http.ResponseWriter, status int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	Json(w, status, errorResponse{Error: msg})
}

func InternalServerError(w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %v", err)
	Err(w, http.StatusInternalServerError, "internal server error")
}

func BadRequest(w http.ResponseWriter, msg string) {
	log.Printf("Bad Request: %v", msg)
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
