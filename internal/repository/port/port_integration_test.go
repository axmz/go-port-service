package port

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    domain "github.com/axmz/go-port-service/internal/domain/port"
    "github.com/axmz/go-port-service/pkg/inmem"
)

func setupIntegrationRepo() *PortRepository {
    db := inmem.New[*Port]()
    return New(db)
}

func TestIntegration_PortRepository_Flow(t *testing.T) {
    repo := setupIntegrationRepo()
    ctx := context.Background()

    var p *domain.Port
    var err error

    t.Run("create and upload port", func(t *testing.T) {
        p, err = domain.New("id42", "Port42", "C42", "City42", "Country42", nil, nil, nil, "", "", nil)
        require.NoError(t, err, "failed to create domain port")
        require.NoError(t, repo.UploadPort(ctx, p), "UploadPort failed")
    })

    t.Run("get by ID", func(t *testing.T) {
        got, err := repo.GetPortByID(ctx, "id42")
        require.NoError(t, err, "GetPortByID failed")
        assert.Equal(t, "id42", got.ID())
        assert.Equal(t, "Port42", got.Name())
    })

    t.Run("get all ports", func(t *testing.T) {
        all, err := repo.GetAllPorts(ctx)
        require.NoError(t, err, "GetAllPorts failed")
        assert.Len(t, all, 1)
    })

    t.Run("get ports count", func(t *testing.T) {
        count := repo.GetPortsCount(ctx)
        assert.Equal(t, 1, count)
    })

    t.Run("delete port", func(t *testing.T) {
        deleted, err := repo.DeletePortByID(ctx, "id42")
        require.NoError(t, err, "DeletePortByID failed")
        assert.Equal(t, "id42", deleted.ID())
    })

    t.Run("ensure port is gone", func(t *testing.T) {
        _, err := repo.GetPortByID(ctx, "id42")
        assert.Error(t, err, "expected error for missing port after delete")
    })
}