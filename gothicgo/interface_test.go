package gothicgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterface(t *testing.T) {
	p, err := NewPackage("test")
	assert.NoError(t, err)
	file := p.File("testFile")
	i, err := file.NewInterface("Stringer")
	assert.NoError(t, err)
	i.AddMethod("String", nil, []Type{StringType}, false)
	s := i.str()
	assert.Contains(t, s, "type Stringer interface{")
	assert.Contains(t, s, "String() string")
}
