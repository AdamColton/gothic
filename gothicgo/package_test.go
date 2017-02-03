package gothicgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackage(t *testing.T) {
	p := NewPackage("testPackage")
	assert.Equal(t, 0, len(p.files))
}

func TestBuildComment(t *testing.T) {
	expected := "// this is a test this is a test FOO this is a test this is a test this is a\n// FOO test this is a test this is a test this FOO is a test this is a test\n// this is a test this is a test this is a test\n\n"
	actual := BuildComment("this is a test this is a test FOO this is a test this is a test this is a FOO test this is a test this is a test this FOO is a test this is a test this is a test this is a test this is a test", 80)
	assert.Equal(t, expected, actual)
}
