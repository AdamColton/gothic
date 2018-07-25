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

	cib, err := f.ConstIotaBlock(Uint64Type,
		"Apple",
		"Bannana",
		"Cantaloup",
		"Date",
		"Elderberry",
	)
	assert.NoError(t, err)

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

func TestConstIotaBlockError(t *testing.T) {
	p, err := NewPackage("constIotaErrorTestPkg")
	assert.NoError(t, err)
	f := p.File("constIotaErrorTestFile")
	f.NewFunc("Apple")

	cib, err := f.ConstIotaBlock(Uint64Type,
		"Apple",
		"Bannana",
		"Cantaloup",
		"Date",
		"Elderberry",
	)
	assert.Error(t, err)
	assert.Nil(t, cib)
}
