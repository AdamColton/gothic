package entbp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackage(t *testing.T) {
	p := NewPackage("test")
	assert.Equal(t, "test", p.Name())
	assert.NotNil(t, p.GoPackage())
}
