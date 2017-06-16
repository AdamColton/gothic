package serialbp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/sai"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeMap(t *testing.T) {
	ctx := New()
	ctx.GetPkg().ImportResolver().Add("serial", "github.com/adamcolton/gothic/serial")

	typ := gothicgo.MapOf(gothicgo.StringType, gothicgo.IntType)
	_, err := ctx.serializeMapFunc(typ)
	assert.NoError(t, err)

	wc := sai.New()
	f := ctx.GetPkg().File("serial.gothic")
	f.Package().ImportResolver()
	f.Writer = wc
	f.Prepare()
	f.Generate()

	expectStrs := []string{
		"github.com/adamcolton/gothic/serial",
		"func MarshalMapstringToint(s map[string]int) []byte",
		"func UnmarshalMapstringToint(b *[]byte) map[string]int",
	}
	got := wc.String()
	for _, str := range expectStrs {
		assert.Contains(t, got, str)
	}
}
