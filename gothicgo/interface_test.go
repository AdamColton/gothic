package gothicgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterface(t *testing.T) {
	f := NewFuncSig("String")
	f.Returns(StringType.AsRet())
	i := NewInterface()
	i.AddMethod(f)

	s := i.String()
	assert.Contains(t, s, "interface {\n")
	assert.Contains(t, s, "\tString() string")
}

func TestEmbed(t *testing.T) {
	f := NewFuncSig("String")
	f.Returns(StringType.AsRet())
	i := NewInterface()
	i.AddMethod(f)

	pkg := MustPackage("testInterfaceEmbed")
	td, err := pkg.NewInterfaceTypeDef("Stringer", i)
	assert.NoError(t, err)

	e := NewInterface()
	e.Embed(td)
	f2 := NewFuncSig("Int")
	f2.Returns(IntType.AsRet())
	e.AddMethod(f2)

	s := e.String()
	assert.Contains(t, s, "interface {\n")
	assert.Contains(t, s, "\ttestInterfaceEmbed.Stringer")
	assert.Contains(t, s, "\tInt() int")
}

func TestEmbedExternal(t *testing.T) {
	i := NewExternalInterfaceType(MustPackageRef("foo/bar"), "Bar")

	e := NewInterface()
	e.Embed(i)
	f2 := NewFuncSig("Int")
	f2.Returns(IntType.AsRet())
	e.AddMethod(f2)

	s := e.String()
	assert.Contains(t, s, "interface {\n")
	assert.Contains(t, s, "\tbar.Bar")
	assert.Contains(t, s, "\tInt() int")
}
