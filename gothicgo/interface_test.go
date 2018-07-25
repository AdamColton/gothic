package gothicgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterface(t *testing.T) {
	i := NewInterface()
	i.AddMethod("String", nil, []Type{StringType}, false)

	s := i.String()
	assert.Contains(t, s, "interface {")
	assert.Contains(t, s, "String() string")
}

// func TestDefInterface(t *testing.T) {
// 	p := MustPackageRef("test")
// 	i := DefInterface(p, "Tester")

// 	assert.Equal(t, i.String(), "test.Tester")
// 	assert.Equal(t, i.File(), (*File)(nil))
// 	assert.Equal(t, i.Kind(), InterfaceKind)
// 	assert.Equal(t, i.PackageRef(), p)
// 	buf := &bytes.Buffer{}
// 	i.PrefixWriteTo(buf, NewImports(p))
// 	assert.Equal(t, "Tester", buf.String())

// 	buf.Reset()
// 	i.PrefixWriteTo(buf, NewImports(MustPackageRef("foo")))
// 	assert.Equal(t, "test.Tester", buf.String())
// }
