package gothicgo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstIotaBlock(t *testing.T) {
	assert.True(t, true)
	cib := ConstIotaBlock{
		Type: Uint64Type,
		Rows: []string{
			"Apple",
			"Bannana",
			"Cantaloup",
			"Date",
			"Elderberry",
		},
	}

	buf := &bytes.Buffer{}
	_, err := cib.WriteTo(buf)
	assert.NoError(t, err)

	expected := `const (
	Apple uint64 = iota
	Bannana
	Cantaloup
	Date
	Elderberry
)
`

	assert.Equal(t, expected, buf.String())
}
