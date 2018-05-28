package gothicgo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackage(t *testing.T) {
	p, err := NewPackage("testPackage")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(p.files))
}

func TestComment(t *testing.T) {
	expected := "// this is a test this is a test FOO this is a test this is a test this is a FOO\n// test this is a test this is a test this FOO is a test this is a test this is\n// a test this is a test this is a test\n"
	var buf bytes.Buffer
	Comment{
		Comment: "this is a test this is a test FOO this is a test this is a test this is a FOO test this is a test this is a test this FOO is a test this is a test this is a test this is a test this is a test",
		Width:   80,
	}.WriteTo(&buf)
	assert.Equal(t, expected, buf.String())
}

func TestPackageRef(t *testing.T) {
	pr, err := NewPackageRef("this should/fail")
	assert.Error(t, err)
	assert.Nil(t, pr)

	pass := "this/should/pass"
	pr, err = NewPackageRef(pass)
	assert.NoError(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, pass, pr.String())
	assert.Equal(t, "pass", pr.Name())
}
