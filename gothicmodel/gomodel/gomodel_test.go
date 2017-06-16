package gomodel

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoModel(t *testing.T) {
	m := gothicmodel.New("test").
		AddField("Name", "string").
		AddField("Age", "int")

	pkg := gothicgo.NewPackage("test")
	gm := Struct(pkg, m)

	gm.Prepare()
	s := gm.String()
	assert.Contains(t, s, "type test struct {")
	assert.Regexp(t, "Name +string", s)
	assert.Regexp(t, "Age +int", s)
}
