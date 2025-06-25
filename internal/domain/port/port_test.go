package port

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	assert.Equal(t, name, p.Name())
	assert.Equal(t, id, p.ID())
}

func TestNewPort_ValidationError(t *testing.T) {
	_, err := New(
		"", "", "CODE", "City", "Country",
		nil, nil, nil, "", "", nil,
	)
	require.Error(t, err, "expected error for missing required fields")
}

func TestSetName(t *testing.T) {
	p, _ := New(
		"ID2", "OldName", "CODE", "City", "Country",
		nil, nil, nil, "", "", nil,
	)
	err := p.SetName("NewName")
	require.NoError(t, err)
	assert.Equal(t, "NewName", p.Name())

	err = p.SetName("")
	require.Error(t, err, "expected error for empty name")
}
