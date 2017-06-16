package gothicmodel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel(t *testing.T) {
	m := New("testModel").
		AddPrimary("id", "[]byte").
		AddField("Name", "string").
		AddField("Age", "int")

	assert.NotNil(t, m)
	assert.Equal(t, "testModel", m.Name())
	assert.Equal(t, "id", m.Primary())
	assert.Equal(t, []string{"id", "Name", "Age"}, m.Fields())
	s, ok := m.Field("id")
	assert.True(t, ok)
	assert.Equal(t, "[]byte", s)
}
