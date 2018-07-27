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
		nameType: NameType{
			N: "bar",
			T: PointerTo(NewExternalType(pkg, "bar")),
		},
		Tags: make(map[string]string),
	}
	file := pkg.File("foo")
	assert.Equal(t, "bar *bar", typeToString(f, file))

	assert.Equal(t, "bar *foo.bar", f.String())

	f.nameType.N = ""
	assert.Equal(t, "*foo.bar", f.String())

	f.Tags["someKey"] = "someValue"
	assert.Equal(t, "*foo.bar `someKey:\"someValue\"`", f.String())

	f.nameType.N = "glorp"
	assert.Equal(t, "glorp *foo.bar `someKey:\"someValue\"`", f.String())

}

func TestFieldMethods(t *testing.T) {
	f := &Field{
		nameType: NameType{
			N: "bar",
			T: PointerTo(NewExternalType(MustPackageRef("foo"), "bar")),
		},
		Tags: map[string]string{"key": "value"},
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
	expected := `struct {
	foo *foo.Foo
	bar *foo.Bar
	*foo.Glorp
}`

	foo := MustPackageRef("foo")

	s := NewStruct()
	s.AddField("foo", PointerTo(NewExternalType(foo, "Foo")))
	s.AddField("bar", PointerTo(NewExternalType(foo, "Bar")))
	s.Embed(PointerTo(NewExternalType(foo, "Glorp")))

	assert.Equal(t, expected, s.String())
}

func TestWriteStruct(t *testing.T) {
	s := NewStruct()
	s.AddField("time", NewExternalType(MustPackageRef("time"), "Time"))

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
