package gomodel

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoModel(t *testing.T) {
	TypeTags["datetime"] = Tags{
		"json": "-",
	}
	m, err := gothicmodel.New("test", gothicmodel.Fields{
		{"Name", "string"},
		{"Age", "int"},
		{"LastLogin", "datetime"},
	})
	assert.NoError(t, err)

	pkg, err := gothicgo.NewPackage("test")
	assert.NoError(t, err)
	gm, err := Struct(pkg, m)
	assert.NoError(t, err)

	gm.Prepare()
	s := gm.String()
	assert.Contains(t, s, "type test struct {")
	assert.Regexp(t, "Name +string", s)
	assert.Regexp(t, "Age +int", s)
	assert.Regexp(t, "LastLogin +time.Time +`json:\"-\"`", s)
	println(s)

	fs := m.Fields()
	expected := []Field{
		Field{
			base: fs[0],
			kind: gothicgo.StringType,
		},
		Field{
			base: fs[1],
			kind: gothicgo.IntType,
		},
		Field{
			base: fs[2],
			kind: Types["datetime"],
		},
	}
	gmfs := gm.Fields()
	assert.Equal(t, expected, gmfs)
	assert.Equal(t, "Name", gmfs[0].Name())
	assert.Equal(t, "string", gmfs[0].Type())
	assert.Equal(t, true, gmfs[0].Primary())
	assert.Equal(t, false, gmfs[1].Primary())
	assert.Equal(t, gothicgo.StringType, gmfs[0].GoType())

	f, ok := m.Field("Age")
	assert.True(t, ok)
	assert.Equal(t, "Age", f.Name())
}
