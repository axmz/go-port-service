package port

import (
    "context"
    "testing"

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
        if err != nil {
            t.Fatalf("failed to create domain port: %v", err)
        }
        if err := repo.UploadPort(ctx, p); err != nil {
            t.Fatalf("UploadPort failed: %v", err)
        }
    })

    t.Run("get by ID", func(t *testing.T) {
        got, err := repo.GetPortByID(ctx, "id42")
        if err != nil {
            t.Fatalf("GetPortByID failed: %v", err)
        }
        if got.ID() != "id42" || got.Name() != "Port42" {
            t.Errorf("unexpected port data: %+v", got)
        }
    })

    t.Run("get all ports", func(t *testing.T) {
        all, err := repo.GetAllPorts(ctx)
        if err != nil {
            t.Fatalf("GetAllPorts failed: %v", err)
        }
        if len(all) != 1 {
            t.Errorf("expected 1 port, got %d", len(all))
        }
    })

    t.Run("get ports count", func(t *testing.T) {
        count := repo.GetPortsCount(ctx)
        if count != 1 {
            t.Errorf("expected count 1, got %d", count)
        }
    })

    t.Run("delete port", func(t *testing.T) {
        deleted, err := repo.DeletePortByID(ctx, "id42")
        if err != nil {
            t.Fatalf("DeletePortByID failed: %v", err)
        }
        if deleted.ID() != "id42" {
            t.Errorf("expected deleted ID 'id42', got %s", deleted.ID())
        }
    })

    t.Run("ensure port is gone", func(t *testing.T) {
        _, err := repo.GetPortByID(ctx, "id42")
        if err == nil {
            t.Errorf("expected error for missing port after delete")
        }
    })
}