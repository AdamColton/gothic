package serialbp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/sai"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializePtr(t *testing.T) {
	typ := gothicgo.PointerTo(gothicgo.IntType)
	serialDef := serializePtrFunc(typ)
	assert.Equal(t, SerialHelperPackage, serialDef.PackageName())

	wc := sai.New()
	f := serialHelperPackage().File("serial.gothic")
	f.Writer = wc
	f.Prepare()
	f.Generate()

	expectStrs := []string{
		"DO NOT MODIFY",
		"package serialHelpers",
		"github.com/adamcolton/gothic/serial",
		"func MarshalPtrint(s *int) []byte",
		"func UnmarshalPtrint(b *[]byte) *int",
	}

	got := wc.String()
	for _, str := range expectStrs {
		assert.Contains(t, got, str)
	}

	// clear serial helper package for other tests
	shp = nil
}
