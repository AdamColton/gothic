package gothicgo

import (
	"github.com/adamcolton/gothic/gothicio"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeDef(t *testing.T) {
	p, err := NewPackage("test")
	assert.NoError(t, err)
	f := p.File("testfile")
	td := f.NewTypeDef("TypeDefTest", IntType)
	var tp Type = td
	assert.NotNil(t, tp)

	td.NewMethod("Test").Body = gothicio.StringWriterTo(`fmt.Println("Hello, world")`)
	f.AddNameImports("fmt")

	s := f.String()

	assert.Contains(t, s, "func (t *TypeDefTest) Test()")
	assert.Contains(t, s, `"fmt"`)
	assert.Contains(t, s, `type TypeDefTest int`)
}
