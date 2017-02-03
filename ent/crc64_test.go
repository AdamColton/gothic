package ent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCRC64(t *testing.T) {
	a := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	b := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	assert.Equal(t, CRC64(a), CRC64(b), "CRC values of equal byte slices should be equal")
}
