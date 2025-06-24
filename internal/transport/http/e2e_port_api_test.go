package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	repository "github.com/axmz/go-port-service/internal/repository/port"
	services "github.com/axmz/go-port-service/internal/services/port"
	handlers "github.com/axmz/go-port-service/internal/transport/http/handlers"
	"github.com/axmz/go-port-service/internal/transport/http/response"
	router "github.com/axmz/go-port-service/internal/transport/http/router"
	"github.com/axmz/go-port-service/pkg/inmem"
)

const (
	portsJsonPath = "../../../ports.json"
	portsCount    = 1632.0
	sampleID      = "ZWUTA"
)

func TestE2E_PortAPI(t *testing.T) {
	// Setup in-memory repo, service, handlers, and router
	db := inmem.New[*repository.Port]()
	repo := repository.New(db)
	svc := services.New(repo)
	h := handlers.NewHTTPHandlers(svc)
	r := router.Router(h, nil)

	// Load ports.json
	portsJson, err := os.ReadFile(portsJsonPath)
	if err != nil {
		t.Fatalf("failed to read ports.json: %v", err)
	}

	t.Run("upload ports", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/ports", bytes.NewReader(portsJson))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		resp := response.Response{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal upload response: %v", err)
		}
		if resp.Status != "OK" || resp.Data != portsCount {
			t.Fatalf("expected status OK and data %f, got status %s and data %d", portsCount, resp.Status, resp.Data)
		}
	})

	t.Run("get port by id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/ports/"+sampleID, nil)
		req.SetPathValue("id", sampleID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 on get by id, got %d", w.Code)
		}
		body, _ := io.ReadAll(w.Body)
		if !bytes.Contains(body, []byte(sampleID)) {
			t.Errorf("expected response to contain port ID %s", sampleID)
		}
	})

	getPortsCount := func(t *testing.T, want float64) {
		req := httptest.NewRequest("GET", "/api/ports/count", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 on count, got %d", w.Code)
		}
		resp := response.Response{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode count: %v", err)
		}
		if resp.Data != want {
			t.Errorf("expected count %f, got %d", want, resp.Data)
		}
	}

	t.Run("get ports count", func(t *testing.T) {
		getPortsCount(t, portsCount)
	})

	t.Run("delete port by id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/ports/"+sampleID, nil)
		req.SetPathValue("id", sampleID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 on delete, got %d", w.Code)
		}
	})

	t.Run("get ports count after delete", func(t *testing.T) {
		getPortsCount(t, portsCount-1)
	})
}
