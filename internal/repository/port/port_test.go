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
	t.Run("UploadAndGet", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		p := testDomainPort()

		require.NoError(t, repo.Upload(ctx, p), "Upload failed")

		got, err := repo.Get(ctx, "id1")
		require.NoError(t, err, "Get failed")
		assert.Equal(t, p.ID(), got.ID(), "expected IDs to match")
	})

	t.Run("Get_NotFound", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()

		_, err := repo.Get(ctx, "notfound")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("GetAll", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		p := testDomainPort()
		_ = repo.Upload(ctx, p)

		ports, err := repo.GetAll(ctx)
		require.NoError(t, err, "GetAll failed")
		assert.Len(t, ports, 1)
	})

	t.Run("Count", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		assert.Equal(t, 0, repo.Count(ctx), "expected count 0")
		_ = repo.Upload(ctx, testDomainPort())
		assert.Equal(t, 1, repo.Count(ctx), "expected count 1")
	})

	t.Run("Delete", func(t *testing.T) {
		t.Parallel()
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		_ = repo.Upload(ctx, testDomainPort())

		p, err := repo.Delete(ctx, "id1")
		require.NoError(t, err, "Delete failed")
		assert.Equal(t, "id1", p.ID(), "expected deleted ID 'id1'")
		_, err = repo.Get(ctx, "id1")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		mem := newMockInMem()
		repo := New(mem)
		ctx := context.Background()
		_, err := repo.Delete(ctx, "notfound")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
