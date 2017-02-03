package gothicgo

import (
	"github.com/adamcolton/sai"
	"github.com/stretchr/testify/assert"
	"testing"
)

type istruct string

func (i istruct) PackageName() string {
	return string(i)
}

func TestFieldString(t *testing.T) {
	f := &Field{
		nameType: &NameType{
			N: "bar",
			T: PointerTo(DefStruct("foo.bar")),
		},
		tags: map[string]string{},
		stct: istruct("foo"),
	}
	assert.Equal(t, "bar *bar", f.String())

	f.stct = istruct("notFoo")
	assert.Equal(t, "bar *foo.bar", f.String())

	f.nameType.N = ""
	assert.Equal(t, "*foo.bar", f.String())

	f.tags["someKey"] = "someValue"
	assert.Equal(t, "*foo.bar `someKey:\"someValue\"`", f.String())

	f.nameType.N = "glorp"
	assert.Equal(t, "glorp *foo.bar `someKey:\"someValue\"`", f.String())

}

func TestFieldMethods(t *testing.T) {
	f := &Field{
		nameType: &NameType{
			N: "bar",
			T: PointerTo(DefStruct("foo.bar")),
		},
		tags: map[string]string{"key": "value"},
		stct: istruct("foo"),
	}

	assert.Equal(t, "bar", f.Name())
	assert.Equal(t, "*foo.bar", f.Type().String())

	tg, ok := f.Tag("key")
	assert.True(t, ok)
	assert.Equal(t, "value", tg)

	tg, ok = f.Tag("key2")
	assert.False(t, ok, "Did not expected 'key2'")

	f.SetTag("key2", "value2")
	tg, ok = f.Tag("key2")
	assert.True(t, ok)
	assert.Equal(t, "value2", tg)
}

func TestStructString(t *testing.T) {
	expected := `type test struct {
	foo *foo.Foo
	bar *foo.Bar
	*foo.Glorp
}`

	p := NewPackage("test")
	p.ImportPath = "test"
	p.OutputPath = "test"
	s := p.NewStruct("test")
	s.AddField("foo", PointerTo(DefStruct("foo.Foo")))
	s.AddField("bar", PointerTo(DefStruct("foo.Bar")))
	s.Embed(PointerTo(DefStruct("foo.Glorp")))

	assert.Equal(t, expected, s.String())
}

func TestImportString(t *testing.T) {
	p := NewPackage("test")
	p.ImportPath = "test"
	p.OutputPath = "test"
	s := p.NewStruct("test")
	s.AddField("time", DefStruct("time.Time"))
	s.Prepare()

	time, ok := s.file.Imports.pkgs["time"]
	assert.True(t, ok)
	assert.Equal(t, "", time)
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
	p := NewPackage("test")
	p.ImportPath = "test"
	p.OutputPath = "test"
	s := p.NewStruct("test")
	s.AddField("time", DefStruct("time.Time"))

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
	p := NewPackage("test")
	p.ImportPath = "test"
	p.OutputPath = "test"
	s := p.NewStruct("test")
	s.AddField("time", DefStruct("time.Time"))

	m := s.NewMethod("foo", Arg("name", StringType))
	m.Body = "fmt.Println(\"Hi\", name)"
	m.AddPackageImport("fmt")

	wc := sai.New()
	f := s.file
	f.Writer = wc
	f.Prepare()
	f.Generate()

	assert.Equal(t, expected, wc.String())
}

func TestDefStruct(t *testing.T) {
	s := DefStruct("foo.Bar")
	assert.Equal(t, "foo", s.PackageName())
	assert.Equal(t, "Bar", s.Name())
}
