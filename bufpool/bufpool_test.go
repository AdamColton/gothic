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
