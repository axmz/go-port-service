package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON writes a JSON response with custom status code.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding JSON: %v", err)
	}
}

func JSONOK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

func JSONError(w http.ResponseWriter, status int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	JSON(w, status, errorResponse{Error: msg})
}

func MissingID(w http.ResponseWriter) {
	JSONError(w, http.StatusBadRequest, "missing required 'id' query parameter")
}

func NotFound(w http.ResponseWriter) {
	JSONError(w, http.StatusNotFound, "resource not found")
}

func Text(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func HTML(w http.ResponseWriter, status int, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(html))
}

func InternalServerError(w http.ResponseWriter, err error) {
	log.Printf("internal server error: %v", err)
	JSONError(w, http.StatusInternalServerError, "internal server error")
}
