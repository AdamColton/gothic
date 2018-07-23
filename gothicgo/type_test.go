package gothicgo

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeString(t *testing.T) {
	expected := "map[foo.bar]*person.Person"

	fb := DefStruct(MustPackageRef("foo"), "bar")
	p := PointerTo(DefStruct(MustPackageRef("person"), "Person"))
	tp := MapOf(fb, p)

	assert.Equal(t, expected, tp.String())
	assert.Equal(t, MapKind, tp.Kind())
	if mtp, ok := tp.(MapType); ok {
		assert.Equal(t, "foo.bar", mtp.Key().String())
	} else {
		t.Error("Could not cast to MapType")
	}
}

func TestPadTest(t *testing.T) {
	if fmt.Sprintf("%-8s|", "test") != "test    |" {
		t.Error("You do not understand padding")
	}
}

func TestFuncTypeString(t *testing.T) {
	expected := "func() string"

	f := NewFunc(NewImports(PkgBuiltin()), "")
	f.Returns(Ret(StringType))
	tp := f.Type()

	assert.Equal(t, expected, tp.String())
	assert.Equal(t, FuncKind, tp.Kind())
}

func TestMapToFuncTypeString(t *testing.T) {
	expected := "map[string]func() string"

	f := NewFunc(NewImports(PkgBuiltin()), "")
	f.Returns(Ret(StringType))
	tp := MapOf(StringType, f.Type())

	assert.Equal(t, expected, tp.String())

	if assert.Equal(t, MapKind, tp.Kind()) {
		assert.Equal(t, "func() string", tp.Elem().String())
	}
}

func TestMapType(t *testing.T) {
	person := MustPackageRef("person")
	foo := MustPackageRef("foo")
	fb := DefStruct(foo, "bar")
	p := PointerTo(DefStruct(person, "Person"))
	tp := MapOf(fb, p)

	buf := &bytes.Buffer{}
	tp.PrefixWriteTo(buf, NewImports(person))
	assert.Equal(t, "map[foo.bar]*Person", buf.String(), "Relative to person")

	buf.Reset()
	tp.PrefixWriteTo(buf, NewImports(foo))
	assert.Equal(t, "map[bar]*person.Person", buf.String())

	buf.Reset()
	tp.PrefixWriteTo(buf, DefaultPrefixer)
	assert.Equal(t, "map[foo.bar]*person.Person", buf.String())
}
