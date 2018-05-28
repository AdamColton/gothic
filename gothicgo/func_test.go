package gothicgo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type writeToString string

func (s writeToString) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(s))
	return int64(n), err
}

type writecloser struct {
	*bytes.Buffer
}

func (wc *writecloser) Close() error {
	return nil
}

func TestNameTypeSliceToString(t *testing.T) {
	s := []NameType{
		Arg("first", StringType),
		Arg("middle", StringType),
		Arg("last", StringType),
		Arg("title", StringType),
	}
	assert.Equal(t, "first, middle, last, title string", nameTypeSliceToString(nil, s, false))

	s = []NameType{
		Arg("first", StringType),
		Arg("age", IntType),
	}
	assert.Equal(t, "first string, age int", nameTypeSliceToString(nil, s, false))
}

func TestFuncString(t *testing.T) {
	f := NewFunc(NewImports(PkgBuiltin()), "Foo", Arg("name", StringType), Arg("age", IntType))
	f.Returns(Ret(PointerTo(DefStruct(PkgBuiltin(), "Person"))))
	f.Body = writeToString("\treturn &Person{\n\t\tName: name,\n\t\tAge: age,\n\t}")

	expected := "func Foo(name string, age int) *Person {\n\treturn &Person{\n\t\tName: name,\n\t\tAge: age,\n\t}\n}\n\n"
	got := f.String()
	assert.Equal(t, expected, got)
}

func TestFuncStringVariadic(t *testing.T) {
	f := NewFunc(NewImports(PkgBuiltin()), "Foo", Arg("name", StringType), Arg("code", IntType))
	f.Returns(Ret(PointerTo(DefStruct(PkgBuiltin(), "Person"))))
	f.Body = writeToString("\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}")
	f.Variadic = true

	expected := "func Foo(name string, code ...int) *Person {\n\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}\n}\n\n"
	got := f.String()
	assert.Equal(t, expected, got)
}

func TestWriteFunc(t *testing.T) {
	p, err := NewPackage("test")
	assert.NoError(t, err)
	p.OutputPath = "test"

	wc := &writecloser{
		Buffer: &bytes.Buffer{},
	}
	f := p.File("testFile")
	f.Writer = wc

	fn := f.NewFunc("Foo", Arg("name", StringType), Arg("code", IntType))
	fn.Returns(Ret(PointerTo(DefStruct(MustPackageRef("test"), "Person"))))
	fn.Body = writeToString("\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}")
	fn.Variadic = true

	f.Prepare()
	f.Generate()

	expected := "// This code was generated from a Gothic Blueprint, DO NOT MODIFY\n\npackage test\n\nfunc Foo(name string, code ...int) *Person {\n\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}\n}\n"
	assert.Equal(t, expected, wc.String())
}

func TestFuncType(t *testing.T) {
	f := NewFunc(NewImports(PkgBuiltin()), "Foo", Arg("name", StringType), Arg("age", IntType))
	p, err := NewPackage("test")
	assert.NoError(t, err)
	f.Returns(Ret(PointerTo(DefStruct(p, "Person"))))
	f.File = p.File("testFile")

	ft := f.Type()

	assert.Equal(t, "func(string, int) *test.Person", ft.Name())
	assert.Equal(t, "func(string, int) *test.Person", ft.String())
}

func TestFuncCall(t *testing.T) {
	p, err := NewPackage("test")
	assert.NoError(t, err)
	file := p.File("test")

	args := []NameType{
		Ret(StringType),
		Ret(StringType),
	}
	fc := FuncCall(p, "myFn", args, nil)
	assert.Equal(t, "myFn(Maggie, Bea)", fc.Call(file, "Maggie", "Bea"))
	assert.Equal(t, fc.Args(), args)

	p, err = NewPackage("foo")
	pre := NewImports(p)
	assert.NoError(t, err)
	assert.Equal(t, "test.myFn(Maggie, Bea)", fc.Call(pre, "Maggie", "Bea"))

	fc = file.NewFunc("Foo", Arg("name", StringType), Arg("age", IntType))
	assert.Equal(t, "Foo(adam, 32)", fc.Call(file, "adam", "32"))
	assert.Equal(t, "test.Foo(adam, 32)", fc.Call(pre, "adam", "32"))
}

func TestFuncComment(t *testing.T) {
	f := NewFunc(NewImports(PkgBuiltin()), "Foo")
	f.Body = writeToString("")
	f.Comment = "is a basic function with no args"

	got := f.String()
	assert.Contains(t, got, "// Foo is a basic function with no args")
}
