package gothicgo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstIotaBlock(t *testing.T) {
	p, err := NewPackage("constIotaTestPkg")
	assert.NoError(t, err)
	f := p.File("constIotaTestFile")

	cib := f.ConstIotaBlock(Uint64Type,
		"Apple",
		"Bannana",
		"Cantaloup",
		"Date",
		"Elderberry",
	)

	buf := &bytes.Buffer{}
	_, err = cib.WriteTo(buf)
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
