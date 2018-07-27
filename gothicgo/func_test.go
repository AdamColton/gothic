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

	str, err := nameTypeSliceToString(nil, s, false)
	assert.NoError(t, err)
	assert.Equal(t, "first, middle, last, title string", str)

	s = []NameType{
		Arg("first", StringType),
		Arg("age", IntType),
	}
	str, err = nameTypeSliceToString(nil, s, false)
	assert.NoError(t, err)
	assert.Equal(t, "first string, age int", str)
}

func TestFuncString(t *testing.T) {
	pkg, err := NewPackage("testpkg")
	assert.NoError(t, err)
	file := pkg.File("testfile")
	buf := &bytes.Buffer{}
	file.Writer = buf

	f, err := file.NewFunc("Foo", Arg("name", StringType), Arg("age", IntType))
	assert.NoError(t, err)
	f.Returns(Ret(PointerTo(NewExternalType(PkgBuiltin(), "Person"))))
	f.BodyString("\treturn &Person{\n\t\tName: name,\n\t\tAge: age,\n\t}")

	file.Prepare()
	file.Generate()
}

func TestFuncStringVariadic(t *testing.T) {
	pkg, err := NewPackage("testpkg")
	assert.NoError(t, err)
	file := pkg.File("testfile")
	buf := &bytes.Buffer{}
	file.Writer = buf

	f, err := file.NewFunc("Foo", Arg("name", StringType), Arg("code", IntType))
	assert.NoError(t, err)

	person := PointerTo(NewExternalType(PkgBuiltin(), "Person"))
	person.PrefixWriteTo(buf, file)

	f.Returns(Ret(person))
	f.BodyString("\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}")
	f.Variadic = true

	file.Prepare()
	file.Generate()

	assert.Contains(t, buf.String(), "func Foo(name string, code ...int) *Person {")
}

func TestWriteFunc(t *testing.T) {
	pkg, err := NewPackage("testpkg")
	assert.NoError(t, err)
	f := pkg.File("testfile")
	buf := &bytes.Buffer{}
	f.Writer = buf

	fn, err := f.NewFunc("Foo", Arg("name", StringType), Arg("code", IntType))
	assert.NoError(t, err)
	fn.Returns(Ret(PointerTo(NewExternalType(MustPackageRef("test"), "Person"))))
	fn.BodyString("\treturn &Person{\n\t\tName: name,\n\t\tCode: code,\n\t}")
	fn.Variadic = true

	f.Prepare()
	f.Generate()

	assert.Contains(t, buf.String(), "func Foo(name string, code ...int) *test.Person {")
}

func TestFuncType(t *testing.T) {
	pkg, err := NewPackage("testpkg")
	assert.NoError(t, err)
	file := pkg.File("testfile")

	f, err := file.NewFunc("Foo", Arg("name", StringType), Arg("age", IntType))
	assert.NoError(t, err)

	p, err := NewPackage("test")
	assert.NoError(t, err)
	f.Returns(Ret(PointerTo(NewExternalType(p, "Person"))))

	assert.Equal(t, "func Foo(string, int) *test.Person", f.Type().String())
}

func TestFuncCall(t *testing.T) {
	p, err := NewPackage("test")
	assert.NoError(t, err)
	file := p.File("test")

	f, err := file.NewFunc("myFn", Arg("dog1", StringType), Arg("dog2", StringType))
	assert.NoError(t, err)

	assert.Equal(t, "myFn(Maggie, Bea)", f.Call(file, "Maggie", "Bea"))

	p, err = NewPackage("foo")
	pre := NewImports(p)
	assert.NoError(t, err)
	assert.Equal(t, "test.myFn(Maggie, Bea)", f.Call(pre, "Maggie", "Bea"))
}

func TestFuncComment(t *testing.T) {
	pkg, err := NewPackage("testpkg")
	assert.NoError(t, err)
	file := pkg.File("testfile")

	f, err := file.NewFunc("Foo")
	assert.NoError(t, err)
	f.Comment = "is a basic function with no args"

	got := f.String()
	assert.Contains(t, got, "// Foo is a basic function with no args")
}

func TestExternalFuncCall(t *testing.T) {
	pkg := MustPackageRef("test")
	f := NewExternalFunc(pkg, "myFn", Arg("dog1", StringType), Arg("dog2", StringType))

	p, err := NewPackage("foo")
	pre := NewImports(p)
	assert.NoError(t, err)
	assert.Equal(t, "test.myFn(Maggie, Bea)", f.Call(pre, "Maggie", "Bea"))
}
