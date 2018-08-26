package gothicgo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

type closerBuf struct{ *bytes.Buffer }

func (c *closerBuf) Close() error {
	return nil
}

type istruct string

func (i istruct) Prefix(ref PackageRef) string {
	return NewImports(MustPackageRef(string(i))).Prefix(ref)
}

func TestFieldString(t *testing.T) {
	pkg := MustPackage("foo")
	f := &Field{
		NameType: NameType{"bar", PointerTo(NewExternalType(pkg, "bar"))},
		Tags:     make(map[string]string),
	}
	file := pkg.File("foo")
	assert.Equal(t, "bar *bar", typeToString(f, file))

	assert.Equal(t, "bar *foo.bar", f.String())

	f.N = ""
	assert.Equal(t, "*foo.bar", f.String())

	f.Tags["someKey"] = "someValue"
	assert.Equal(t, "*foo.bar `someKey:\"someValue\"`", f.String())

	f.N = "glorp"
	assert.Equal(t, "glorp *foo.bar `someKey:\"someValue\"`", f.String())

}

func TestFieldMethods(t *testing.T) {
	f := &Field{
		NameType: NameType{"bar", PointerTo(NewExternalType(MustPackageRef("foo"), "bar"))},
		Tags:     map[string]string{"key": "value"},
	}

	assert.Equal(t, "bar", f.Name())
	assert.Equal(t, "*foo.bar", f.Type().String())

	tg, ok := f.Tags["key"]
	assert.True(t, ok)
	assert.Equal(t, "value", tg)

	tg, ok = f.Tags["key2"]
	assert.False(t, ok, "Did not expected 'key2'")

	f.Tags["key2"] = "value2"
	tg, ok = f.Tags["key2"]
	assert.True(t, ok)
	assert.Equal(t, "value2", tg)
}

func TestStructString(t *testing.T) {
	foo := MustPackageRef("foo")

	s := NewStruct(
		PointerTo(NewExternalType(foo, "Foo")).Named("foo"),
		PointerTo(NewExternalType(foo, "Bar")).Named("bar"),
		PointerTo(NewExternalType(foo, "Glorp")),
		StringType,
	)

	str := s.String()
	assert.Contains(t, str, "struct {")
	assert.Contains(t, str, "\tfoo *foo.Foo")
	assert.Contains(t, str, "\tbar *foo.Bar")
	assert.Contains(t, str, "\t*foo.Glorp")
	assert.Contains(t, str, "\tstring\n")
}

func TestWriteStruct(t *testing.T) {
	s := NewStruct(
		NewExternalType(MustPackageRef("time"), "Time").Named("time"),
	)

	buf := &bytes.Buffer{}
	_, err := s.PrefixWriteTo(buf, DefaultPrefixer)
	assert.NoError(t, err)

	assert.Contains(t, buf.String(), "struct {")
	assert.Contains(t, buf.String(), "time time.Time")
	assert.Contains(t, buf.String(), "}")
}

func TestNewExternalType(t *testing.T) {
	s := NewExternalType(MustPackageRef("foo"), "Bar")
	assert.Equal(t, "foo", s.PackageRef().String())
}

func TestStructInStruct(t *testing.T) {
	s1 := NewStruct(
		StringType.Named("Name"),
		IntType.Named("Age"),
	)

	s2 := NewStruct(
		StringType.Named("Role"),
		NameType{"Person", s1},
	)

	s := s2.String()
	assert.Contains(t, s, "\tRole string")
	assert.Contains(t, s, "\tPerson struct {")
	assert.Contains(t, s, "\t\tName string")
	assert.Contains(t, s, "\t\tAge int")
}
