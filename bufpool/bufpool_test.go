package bufpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBufferPool(t *testing.T) {
	buf := Get()
	buf.WriteString("this is a test")
	assert.Equal(t, "this is a test", PutStr(buf))
}

func TestPutStr(t *testing.T) {
	buf := Get()
	buf.WriteString("Hello")
	s := PutStr(buf)
	assert.Equal(t, "Hello", s)
	// another process gets the same buffer
	buf.WriteString("Goodbye")
	assert.Equal(t, "Hello", s)
}
