package gothicgo

import (
	"github.com/adamcolton/sai"
	"github.com/stretchr/testify/assert"
	"testing"
)

type istruct string

func (i istruct) Prefix(ref PackageRef) string {
	return NewImports(MustPackageRef(string(i))).Prefix(ref)
}

func TestFieldString(t *testing.T) {
	f := &Field{
		nameType: &NameType{
			N: "bar",
			T: PointerTo(DefStruct(MustPackageRef("foo"), "bar")),
		},
		Tags: make(map[string]string),
		stct: istruct("foo"),
	}
	assert.Equal(t, "bar *bar", f.String())

	f.stct = istruct("notFoo")
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
		nameType: &NameType{
			N: "bar",
			T: PointerTo(DefStruct(MustPackageRef("foo"), "bar")),
		},
		Tags: map[string]string{"key": "value"},
		stct: istruct("foo"),
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
	expected := `type test struct {
	foo *foo.Foo
	bar *foo.Bar
	*foo.Glorp
}`

	foo := MustPackageRef("foo")
	SetImportPath("")
	p, err := NewPackage("test")
	assert.NoError(t, err)
	p.OutputPath = "test"
	s, err := p.NewStruct("test")
	assert.NoError(t, err)
	s.AddField("foo", PointerTo(DefStruct(foo, "Foo")))
	s.AddField("bar", PointerTo(DefStruct(foo, "Bar")))
	s.Embed(PointerTo(DefStruct(foo, "Glorp")))

	assert.Equal(t, expected, s.String())
}

func TestImportString(t *testing.T) {
	SetImportPath("")
	p, err := NewPackage("test")
	assert.NoError(t, err)
	p.OutputPath = "test"
	s, err := p.NewStruct("test")
	assert.NoError(t, err)
	timePkg := MustPackageRef("time")
	s.AddField("time", DefStruct(timePkg, "Time"))
	s.Prepare()

	time, ok := s.file.Imports.refs[timePkg]
	assert.True(t, ok)
	assert.Equal(t, "", time)
	assert.Equal(t, "time", s.file.GetRefName(timePkg))
}

func TestWriteStruct(t *testing.T) {
	expected := `// This code was generated from a Gothic Blueprint, DO NOT MODIFY

package test

import (
	"time"
)

type test struct {
	time time.Time
}
`
	SetImportPath("")
	p, err := NewPackage("test")
	assert.NoError(t, err)
	p.OutputPath = "test"
	s, err := p.NewStruct("test")
	assert.NoError(t, err)
	s.AddField("time", DefStruct(MustPackageRef("time"), "Time"))

	wc := sai.New()
	f := s.file
	f.Writer = wc
	f.Prepare()
	f.Generate()

	assert.Equal(t, expected, wc.String())
}

func TestMethod(t *testing.T) {
	expected := `// This code was generated from a Gothic Blueprint, DO NOT MODIFY

package test

import (
	"fmt"
	"time"
)

type test struct {
	time time.Time
}

func (t *test) foo(name string) {
	fmt.Println("Hi", name)
}
`
	SetImportPath("")
	p, err := NewPackage("test")
	assert.NoError(t, err)
	p.OutputPath = "test"
	s, err := p.NewStruct("test")
	assert.NoError(t, err)
	s.AddField("time", DefStruct(MustPackageRef("time"), "Time"))

	m := s.NewMethod("foo", Arg("name", StringType))
	m.Body = func() (string, error) { return "fmt.Println(\"Hi\", name)", nil }
	m.AddRefImports(MustPackageRef("fmt"))

	wc := sai.New()
	f := s.file
	f.Writer = wc
	f.Prepare()
	f.Generate()

	assert.Equal(t, expected, wc.String())
}

func TestDefStruct(t *testing.T) {
	s := DefStruct(MustPackageRef("foo"), "Bar")
	assert.Equal(t, "foo", s.PackageRef().String())
	assert.Equal(t, "Bar", s.Name())
}
