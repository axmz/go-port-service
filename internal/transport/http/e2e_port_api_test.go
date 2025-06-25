package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repository "github.com/axmz/go-port-service/internal/repository/port"
	services "github.com/axmz/go-port-service/internal/services/port"
	"github.com/axmz/go-port-service/internal/transport/http/handlers"
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
	h := handlers.New(svc)
	r := router.Router(h, nil)

	// Load ports.json
	portsJson, err := os.ReadFile(portsJsonPath)
	require.NoError(t, err, "failed to read ports.json")

	t.Run("upload ports", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/ports", bytes.NewReader(portsJson))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		resp := response.Response{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err, "failed to unmarshal upload response")
		assert.Equal(t, "OK", resp.Status)
		assert.Equal(t, portsCount, resp.Data)
	})

	t.Run("get port by id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/ports/"+sampleID, nil)
		req.SetPathValue("id", sampleID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "expected 200 on get by id")
		body, _ := io.ReadAll(w.Body)
		assert.Contains(t, string(body), sampleID, "expected response to contain port ID")
	})

	Count := func(t *testing.T, want float64) {
		req := httptest.NewRequest("GET", "/api/ports/count", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "expected 200 on count")
		resp := response.Response{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err, "failed to decode count")
		assert.Equal(t, want, resp.Data)
	}

	t.Run("get ports count", func(t *testing.T) {
		Count(t, portsCount)
	})

	t.Run("delete port by id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/ports/"+sampleID, nil)
		req.SetPathValue("id", sampleID)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "expected 200 on delete")
	})

	t.Run("get ports count after delete", func(t *testing.T) {
		Count(t, portsCount-1)
	})
}
