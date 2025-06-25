package port

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/axmz/go-port-service/internal/domain/port"
	"github.com/stretchr/testify/assert"
)

type mockPortService struct {
	GetFunc    func(ctx context.Context, id string) (*port.Port, error)
	DeleteFunc func(ctx context.Context, id string) (*port.Port, error)
	GetAllFunc func(ctx context.Context) ([]*port.Port, error)
	CountFunc  func(ctx context.Context) int
	UploadFunc func(ctx context.Context, p *port.Port) error
}

func (m *mockPortService) Get(ctx context.Context, id string) (*port.Port, error) {
	return m.GetFunc(ctx, id)
}
func (m *mockPortService) Delete(ctx context.Context, id string) (*port.Port, error) {
	return m.DeleteFunc(ctx, id)
}
func (m *mockPortService) GetAll(ctx context.Context) ([]*port.Port, error) {
	return m.GetAllFunc(ctx)
}
func (m *mockPortService) Count(ctx context.Context) int {
	return m.CountFunc(ctx)
}
func (m *mockPortService) Upload(ctx context.Context, p *port.Port) error {
	return m.UploadFunc(ctx, p)
}

func TestUpload_Success(t *testing.T) {
	mockSvc := &mockPortService{
		UploadFunc: func(ctx context.Context, p *port.Port) error {
			return nil
		},
	}
	h := New(mockSvc)

	body := `{"id1": {"name": "Port1", "city": "City1", "country": "Country1", "code": "C1", "alias": [], "regions": [], "coordinates": [], "province": "", "timezone": "", "unlocs": []}}`
	req := httptest.NewRequest("POST", "/api/ports", strings.NewReader(body))
	w := httptest.NewRecorder()

	h.Upload(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpload_BadJSON(t *testing.T) {
	mockSvc := &mockPortService{}
	h := New(mockSvc)

	req := httptest.NewRequest("POST", "/api/ports", strings.NewReader("notjson"))
	w := httptest.NewRecorder()

	h.Upload(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAll(t *testing.T) {
	h := New(&mockPortService{
		GetAllFunc: func(ctx context.Context) ([]*port.Port, error) {
			return []*port.Port{}, nil
		},
	})
	req := httptest.NewRequest("GET", "/api/ports", nil)
	w := httptest.NewRecorder()
	h.GetAll(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAll_Error(t *testing.T) {
	h := New(&mockPortService{
		GetAllFunc: func(ctx context.Context) ([]*port.Port, error) {
			return nil, errors.New("fail")
		},
	})
	req := httptest.NewRequest("GET", "/api/ports", nil)
	w := httptest.NewRecorder()
	h.GetAll(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGet(t *testing.T) {
	h := New(&mockPortService{
		GetFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return &port.Port{}, nil
		},
	})
	req := httptest.NewRequest("GET", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.Get(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGet_NotFound(t *testing.T) {
	h := New(&mockPortService{
		GetFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return nil, port.ErrNotFound
		},
	})
	req := httptest.NewRequest("GET", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.Get(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGet_Error(t *testing.T) {
	h := New(&mockPortService{
		GetFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return nil, errors.New("fail")
		},
	})
	req := httptest.NewRequest("GET", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.Get(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCount(t *testing.T) {
	h := New(&mockPortService{
		CountFunc: func(ctx context.Context) int {
			return 42
		},
	})
	req := httptest.NewRequest("GET", "/api/ports/count", nil)
	w := httptest.NewRecorder()
	h.Count(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDelete(t *testing.T) {
	h := New(&mockPortService{
		DeleteFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return &port.Port{}, nil
		},
	})
	req := httptest.NewRequest("DELETE", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.Delete(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDelete_NotFound(t *testing.T) {
	h := New(&mockPortService{
		DeleteFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return nil, port.ErrNotFound
		},
	})
	req := httptest.NewRequest("DELETE", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.Delete(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDelete_Error(t *testing.T) {
	h := New(&mockPortService{
		DeleteFunc: func(ctx context.Context, id string) (*port.Port, error) {
			return nil, errors.New("fail")
		},
	})
	req := httptest.NewRequest("DELETE", "/api/ports/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	h.Delete(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
