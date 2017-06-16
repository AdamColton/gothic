package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	Environments("prod", "dev", "test")

	SetString("test").
		All("foo").
		On("dev", "bar")

	assert.Equal(t, "foo", GetString("test"))
	SetEnvironment("dev")
	assert.Equal(t, "bar", GetString("test"))
}
