package gothicgo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestTypeDef(t *testing.T) {
	p, err := NewPackage("test")
	assert.NoError(t, err)
	f := p.File("testfile")
	td, err := f.NewTypeDef("TypeDefTest", IntType)
	assert.NoError(t, err)
	var tp Type = td
	assert.NotNil(t, tp)

	m, err := td.NewMethod("Test")
	assert.NoError(t, err)
	m.BodyString(`fmt.Println("Hello, world")`)
	f.AddNameImports("fmt")

	s := f.String()

	assert.Contains(t, s, "func (t *TypeDefTest) Test()")
	assert.Contains(t, s, `"fmt"`)
	assert.Contains(t, s, `type TypeDefTest int`)
}

func TestImportString(t *testing.T) {
	p, err := NewPackage("test")
	assert.NoError(t, err)

	s := NewStruct()
	timePkg := MustPackageRef("time")
	s.AddField("time", NewExternalType(timePkg, "Time"))

	td, err := p.NewTypeDef("TestTypeDefImports", s)
	assert.NoError(t, err)

	assert.NoError(t, td.file.Prepare())

	time, ok := td.File().Imports.refs[timePkg.String()]
	assert.True(t, ok)
	assert.Equal(t, "", time)
	assert.Equal(t, "time", td.File().GetRefName(timePkg))
}

type testMethodBody struct{}

func (b testMethodBody) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	bt := []byte(`fmt.Println("Hi", name)`)
	w.Write(bt)
	return int64(len(bt)), nil
}

func (b testMethodBody) RegisterImports(i *Imports) {
	i.AddNameImports("fmt")
}

func TestMethod(t *testing.T) {
	s := NewStruct()
	s.AddField("time", NewExternalType(MustPackageRef("time"), "Time"))

	p, err := NewPackage("test")
	assert.NoError(t, err)
	td, err := p.NewTypeDef("Test", s)
	assert.NoError(t, err)

	m, err := td.NewMethod("foo", Arg("name", StringType))
	assert.NoError(t, err)

	m.Body = testMethodBody{}
	m.Comment = "says Hi"

	buf := &bytes.Buffer{}
	f := td.file
	f.Writer = buf
	f.Prepare()
	f.Generate()

	assert.Contains(t, buf.String(), "type Test struct {")
	assert.Contains(t, buf.String(), "// foo says Hi")
	assert.Contains(t, buf.String(), "func (t *Test) foo(name string) {")
	assert.Contains(t, buf.String(), `fmt.Println("Hi", name)`)
}
