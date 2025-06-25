package port

import (
	"context"
	"testing"

	domain "github.com/axmz/go-port-service/internal/domain/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestPortRepository(t *testing.T) {
	t.Run("UploadAndGetPortByID", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		p := testDomainPort()

		require.NoError(t, repo.UploadPort(ctx, p), "UploadPort failed")

		got, err := repo.GetPortByID(ctx, "id1")
		require.NoError(t, err, "GetPortByID failed")
		assert.Equal(t, p.ID(), got.ID(), "expected IDs to match")
	})

	t.Run("GetPortByID_NotFound", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()

		_, err := repo.GetPortByID(ctx, "notfound")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("GetAllPorts", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		p := testDomainPort()
		_ = repo.UploadPort(ctx, p)

		ports, err := repo.GetAllPorts(ctx)
		require.NoError(t, err, "GetAllPorts failed")
		assert.Len(t, ports, 1)
	})

	t.Run("GetPortsCount", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		assert.Equal(t, 0, repo.GetPortsCount(ctx), "expected count 0")
		_ = repo.UploadPort(ctx, testDomainPort())
		assert.Equal(t, 1, repo.GetPortsCount(ctx), "expected count 1")
	})

	t.Run("DeletePortByID", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		_ = repo.UploadPort(ctx, testDomainPort())

		p, err := repo.DeletePortByID(ctx, "id1")
		require.NoError(t, err, "DeletePortByID failed")
		assert.Equal(t, "id1", p.ID(), "expected deleted ID 'id1'")
		_, err = repo.GetPortByID(ctx, "id1")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("DeletePortByID_NotFound", func(t *testing.T) {
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		_, err := repo.DeletePortByID(ctx, "notfound")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
