package gothicgo

import (
	"github.com/adamcolton/sai"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNameTypeSliceToString(t *testing.T) {
	s := []*NameType{Arg("first", StringType), Arg("middle", StringType), Arg("last", StringType), Arg("title", StringType)}
	assert.Equal(t, "first, middle, last, title string", nameTypeSliceToString(s, "", false))

	s = []*NameType{Arg("first", StringType), Arg("age", IntType)}
	assert.Equal(t, "first string, age int", nameTypeSliceToString(s, "", false))
}

func TestFuncString(t *testing.T) {
	f := NewFunc("Foo", Arg("name", StringType), Arg("age", IntType))
	f.Returns(Ret(PointerTo(DefStruct("Person"))))
	f.Body = "\treturn &Person{\n\t\tName: name,\n\t\tAge: age,\n\t}"

	expected := "func Foo(name string, age int) *Person {\n\treturn &Person{\n\t\tName: name,\n\t\tAge: age,\n\t}\n}\n\n"
	got := f.String()
	assert.Equal(t, expected, got)
}

func TestFuncStringVariadic(t *testing.T) {
	f := NewFunc("Foo", Arg("name", StringType), Arg("code", IntType))
	f.Returns(Ret(PointerTo(DefStruct("Person"))))
	f.Body = "\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}"
	f.Variadic = true

	expected := "func Foo(name string, code ...int) *Person {\n\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}\n}\n\n"
	got := f.String()
	assert.Equal(t, expected, got)
}

func TestWriteFunc(t *testing.T) {
	p := NewPackage("test")
	p.ImportPath = "test"
	p.OutputPath = "test"

	wc := sai.New()
	f := p.File("testFile")
	f.Writer = wc

	fn := f.NewFunc("Foo", Arg("name", StringType), Arg("code", IntType))
	fn.Returns(Ret(PointerTo(DefStruct("test.Person"))))
	fn.Body = "\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}"
	fn.Variadic = true

	f.Prepare()
	f.Generate()

	expected := "// This code was generated from a Gothic Blueprint, DO NOT MODIFY\n\npackage test\n\nfunc Foo(name string, code ...int) *Person {\n\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}\n}\n"
	assert.Equal(t, expected, wc.String())
}

func TestFuncType(t *testing.T) {
	f := NewFunc("Foo", Arg("name", StringType), Arg("age", IntType))
	f.Returns(Ret(PointerTo(DefStruct("test.Person"))))
	f.File = NewPackage("test").File("testFile")

	ft := f.Type()

	assert.Equal(t, "func(string, int) *Person", ft.Name())
	assert.Equal(t, "func(string, int) *test.Person", ft.String())
}

func TestFuncSignature(t *testing.T) {
	f := NewFunc("Foo", Arg("name", StringType), Arg("age", IntType))
	f.Returns(Ret(PointerTo(DefStruct("test.Person"))))
	f.File = NewPackage("test").File("testFile")

	fs := f.RelSignature("test")

	assert.Equal(t, "Foo(name string, age int) *Person", fs)
}
