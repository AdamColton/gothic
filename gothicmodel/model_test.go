package gothicmodel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel(t *testing.T) {
	m, err := New("testModel", Fields{
		{"id", "[]byte"},
		{"Name", "string"},
		{"Age", "int"},
	})

	assert.NotNil(t, m)
	assert.NoError(t, err)
	assert.Equal(t, "testModel", m.Name())
	assert.Equal(t, "id", m.Primary().Name())

	f := m.Fields()
	if !assert.Len(t, f, 3) {
		return
	}
	assert.Equal(t, "id", f[0].Name())
	assert.Equal(t, "Name", f[1].Name())
	assert.Equal(t, "Age", f[2].Name())

	s, ok := m.Field("id")
	assert.True(t, ok)
	assert.Equal(t, "[]byte", s.Type())
}

func TestAddPrimary(t *testing.T) {
	m, err := New("primaryTestModel", Fields{
		{"initialPrimary", "[]byte"},
		{"neverPrimary", "string"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "initialPrimary", m.Primary().Name())
	m.AddPrimary("newPrimary", "int")
	assert.Equal(t, "newPrimary", m.Primary().Name())

	fs := m.Fields()
	if assert.Len(t, fs, 3) {
		assert.True(t, fs[0].Primary())
		assert.False(t, fs[1].Primary())
		assert.False(t, fs[2].Primary())

		assert.Equal(t, "newPrimary", fs[0].Name())
		assert.Equal(t, "initialPrimary", fs[1].Name())
		assert.Equal(t, "neverPrimary", fs[2].Name())
	}
}
