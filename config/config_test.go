package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	Environments("prod", "dev", "test")

	SetString("test").
		As("bar", "dev").
		As("foo", "prod", "test")

	assert.Equal(t, "foo", GetString("test"))
	SetEnvironment("dev")
	assert.Equal(t, "bar", GetString("test"))
}
