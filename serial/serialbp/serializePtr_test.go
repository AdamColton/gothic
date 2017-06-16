package serialbp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/sai"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializePtr(t *testing.T) {
	ctx := New()
	typ := gothicgo.PointerTo(gothicgo.IntType)
	_, err := ctx.serializePtrFunc(typ)
	assert.NoError(t, err)

	wc := sai.New()
	f := ctx.GetPkg().File("serial.gothic")
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
}
