package gothicgo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterface(t *testing.T) {
	p, err := NewPackage("test")
	assert.NoError(t, err)
	file := p.File("testFile")
	i, err := file.NewInterface("Stringer")
	assert.NoError(t, err)
	i.AddMethod("String", nil, []Type{StringType}, false)
	var buf bytes.Buffer
	i.WriteTo(&buf)
	s := buf.String()
	assert.Contains(t, s, "type Stringer interface{")
	assert.Contains(t, s, "String() string")

	var tp Type = i
	assert.Equal(t, "test.Stringer", tp.String())
}

func TestDefInterface(t *testing.T) {
	p := MustPackageRef("test")
	i := DefInterface(p, "Tester")

	assert.Equal(t, i.Name(), "Tester")
	assert.Equal(t, i.String(), "test.Tester")
	assert.Equal(t, i.File(), (*File)(nil))
	assert.Equal(t, i.Kind(), InterfaceKind)
	assert.Equal(t, i.PackageRef(), p)
	assert.Equal(t, i.RelStr(NewImports(p)), "Tester")
	assert.Equal(t, i.RelStr(NewImports(MustPackageRef("foo"))), "test.Tester")
}
