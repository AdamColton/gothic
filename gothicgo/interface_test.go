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
	assert.Contains(t, s, "interface {")
	assert.Contains(t, s, "String() string")
}
