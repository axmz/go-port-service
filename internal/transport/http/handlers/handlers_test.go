package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/axmz/go-port-service/internal/domain/port"
)

type mockPortService struct {
	GetPortByIDFunc    func(ctx context.Context, id string) (*port.Port, error)
	DeletePortByIDFunc func(ctx context.Context, id string) (*port.Port, error)
	GetAllPortsFunc    func(ctx context.Context) ([]*port.Port, error)
	GetPortsCountFunc  func(ctx context.Context) int
	UploadPortFunc     func(ctx context.Context, p *port.Port) error
}

func (m *mockPortService) GetPortByID(ctx context.Context, id string) (*port.Port, error) {
	return m.GetPortByIDFunc(ctx, id)
}
func (m *mockPortService) DeletePortByID(ctx context.Context, id string) (*port.Port, error) {
	return m.DeletePortByIDFunc(ctx, id)
}
func (m *mockPortService) GetAllPorts(ctx context.Context) ([]*port.Port, error) {
	return m.GetAllPortsFunc(ctx)
}
func (m *mockPortService) GetPortsCount(ctx context.Context) int {
	return m.GetPortsCountFunc(ctx)
}
func (m *mockPortService) UploadPort(ctx context.Context, p *port.Port) error {
	return m.UploadPortFunc(ctx, p)
}

func TestUploadPorts_Success(t *testing.T) {
	mockSvc := &mockPortService{
		UploadPortFunc: func(ctx context.Context, p *port.Port) error {
			return nil
		},
	}
	h := NewHTTPHandlers(mockSvc)

	body := `{"id1": {"name": "Port1", "city": "City1", "country": "Country1", "code": "C1", "alias": [], "regions": [], "coordinates": [], "province": "", "timezone": "", "unlocs": []}}`
	req := httptest.NewRequest("POST", "/api/ports", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.UploadPorts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUploadPorts_BadJSON(t *testing.T) {
	mockSvc := &mockPortService{}
	h := NewHTTPHandlers(mockSvc)

	req := httptest.NewRequest("POST", "/api/ports", strings.NewReader("notjson"))
	w := httptest.NewRecorder()

	h.UploadPorts(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetAllPorts(t *testing.T) {
	h := NewHTTPHandlers(&mockPortService{
		GetAllPortsFunc: func(ctx context.Context) ([]*port.Port, error) {
			return []*port.Port{}, nil
		},
	})
	req := httptest.NewRequest("GET", "/api/ports", nil)
	w := httptest.NewRecorder()
	h.GetAllPorts(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetAllPorts_Error(t *testing.T) {
	h := NewHTTPHandlers(&mockPortService{
		GetAllPortsFunc: func(ctx context.Context) ([]*port.Port, error) {
			return nil, errors.New("fail")
		},
	})
	req := httptest.NewRequest("GET", "/api/ports", nil)
	w := httptest.NewRecorder()
	h.GetAllPorts(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetPortByID(t *testing.T) {
	h := NewHTTPHandlers(&mockPortService{
		GetPortByIDFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return &port.Port{}, nil
		},
	})
	req := httptest.NewRequest("GET", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.GetPortByID(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetPortByID_NotFound(t *testing.T) {
	h := NewHTTPHandlers(&mockPortService{
		GetPortByIDFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return nil, port.ErrNotFound
		},
	})
	req := httptest.NewRequest("GET", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.GetPortByID(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetPortByID_Error(t *testing.T) {
	h := NewHTTPHandlers(&mockPortService{
		GetPortByIDFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return nil, errors.New("fail")
		},
	})
	req := httptest.NewRequest("GET", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.GetPortByID(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetPortsCount(t *testing.T) {
	h := NewHTTPHandlers(&mockPortService{
		GetPortsCountFunc: func(ctx context.Context) int {
			return 42
		},
	})
	req := httptest.NewRequest("GET", "/api/ports/count", nil)
	w := httptest.NewRecorder()
	h.GetPortsCount(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDeletePortByID(t *testing.T) {
	h := NewHTTPHandlers(&mockPortService{
		DeletePortByIDFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return &port.Port{}, nil
		},
	})
	req := httptest.NewRequest("DELETE", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.DeletePortByID(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDeletePortByID_NotFound(t *testing.T) {
	h := NewHTTPHandlers(&mockPortService{
		DeletePortByIDFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return nil, port.ErrNotFound
		},
	})
	req := httptest.NewRequest("DELETE", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.DeletePortByID(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeletePortByID_Error(t *testing.T) {
	h := NewHTTPHandlers(&mockPortService{
		DeletePortByIDFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return nil, errors.New("fail")
		},
	})
	req := httptest.NewRequest("DELETE", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.DeletePortByID(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
