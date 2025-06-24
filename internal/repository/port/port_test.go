package port

import (
	"context"
	"errors"
	"testing"

	domain "github.com/axmz/go-port-service/internal/domain/port"
)

type mockInMem struct {
	store map[string]*Port
}

func newMockInMem() *mockInMem {
	return &mockInMem{store: make(map[string]*Port)}
}

func (m *mockInMem) Get(ctx context.Context, key string) (*Port, bool) {
	val, ok := m.store[key]
	return val, ok
}
func (m *mockInMem) GetAll(ctx context.Context) []*Port {
	res := make([]*Port, 0, len(m.store))
	for _, v := range m.store {
		res = append(res, v)
	}
	return res
}
func (m *mockInMem) Put(ctx context.Context, key string, value *Port) {
	m.store[key] = value
}
func (m *mockInMem) Delete(ctx context.Context, key string) (*Port, bool) {
	val, ok := m.store[key]
	if ok {
		delete(m.store, key)
	}
	return val, ok
}
func (m *mockInMem) Len(ctx context.Context) int {
	return len(m.store)
}

// helpers for conversion
func testDomainPort() *domain.Port {
	p, _ := domain.New("id1", "name", "code", "city", "country", nil, nil, nil, "", "", nil)
	return p
}
func testRepoPort() *Port {
	return &Port{
		ID:      "id1",
		Name:    "name",
		Code:    "code",
		City:    "city",
		Country: "country",
	}
}

func TestPortRepository_UploadAndGetPortByID(t *testing.T) {
	t.Parallel()
	mem := newMockInMem()
	repo := New(mem)
	ctx := context.Background()
	p := testDomainPort()

	// Upload
	if err := repo.UploadPort(ctx, p); err != nil {
		t.Fatalf("UploadPort failed: %v", err)
	}

	// Get
	got, err := repo.GetPortByID(ctx, "id1")
	if err != nil {
		t.Fatalf("GetPortByID failed: %v", err)
	}
	if got.ID() != p.ID() {
		t.Errorf("expected ID %s, got %s", p.ID(), got.ID())
	}
}

func TestPortRepository_GetPortByID_NotFound(t *testing.T) {
	t.Parallel()
	mem := newMockInMem()
	repo := New(mem)
	ctx := context.Background()

	_, err := repo.GetPortByID(ctx, "notfound")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPortRepository_GetAllPorts(t *testing.T) {
	t.Parallel()
	mem := newMockInMem()
	repo := New(mem)
	ctx := context.Background()
	p := testDomainPort()
	_ = repo.UploadPort(ctx, p)

	ports, err := repo.GetAllPorts(ctx)
	if err != nil {
		t.Fatalf("GetAllPorts failed: %v", err)
	}
	if len(ports) != 1 {
		t.Errorf("expected 1 port, got %d", len(ports))
	}
}

func TestPortRepository_GetPortsCount(t *testing.T) {
	t.Parallel()
	mem := newMockInMem()
	repo := New(mem)
	ctx := context.Background()
	if repo.GetPortsCount(ctx) != 0 {
		t.Errorf("expected count 0")
	}
	_ = repo.UploadPort(ctx, testDomainPort())
	if repo.GetPortsCount(ctx) != 1 {
		t.Errorf("expected count 1")
	}
}

func TestPortRepository_DeletePortByID(t *testing.T) {
	t.Parallel()
	mem := newMockInMem()
	repo := New(mem)
	ctx := context.Background()
	_ = repo.UploadPort(ctx, testDomainPort())

	p, err := repo.DeletePortByID(ctx, "id1")
	if err != nil {
		t.Fatalf("DeletePortByID failed: %v", err)
	}
	if p.ID() != "id1" {
		t.Errorf("expected deleted ID 'id1', got %s", p.ID())
	}
	_, err = repo.GetPortByID(ctx, "id1")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestPortRepository_DeletePortByID_NotFound(t *testing.T) {
	mem := newMockInMem()
	repo := New(mem)
	ctx := context.Background()
	_, err := repo.DeletePortByID(ctx, "notfound")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
