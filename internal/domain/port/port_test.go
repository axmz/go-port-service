package port

import (
	"testing"
)

func TestNewPort_Success(t *testing.T) {
	const id = "ID1"
	const name = "PortName"
	p, err := New(
		id, name, "CODE", "City", "Country",
		[]string{"alias1"}, []string{"region1"},
		[]float64{1.23, 4.56}, "Province", "Timezone",
		[]string{"UNLOC1"},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Name() != name {
		t.Errorf("expected name 'PortName', got %s", p.Name())
	}
	if p.ID() != id {
		t.Errorf("expected id 'ID1', got %s", p.ID())
	}
}

func TestNewPort_ValidationError(t *testing.T) {
	_, err := New(
		"", "", "CODE", "City", "Country",
		nil, nil, nil, "", "", nil,
	)
	if err == nil {
		t.Fatal("expected error for missing required fields, got nil")
	}
}

func TestSetName(t *testing.T) {
	p, _ := New(
		"ID2", "OldName", "CODE", "City", "Country",
		nil, nil, nil, "", "", nil,
	)
	err := p.SetName("NewName")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Name() != "NewName" {
		t.Errorf("expected name 'NewName', got %s", p.Name())
	}
	err = p.SetName("")
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}
