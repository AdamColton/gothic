package serialbp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/sai"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeMap(t *testing.T) {
	serialHelperPackage().ImportResolver().Add("serial", "github.com/adamcolton/gothic/serial")

	typ := gothicgo.MapOf(gothicgo.StringType, gothicgo.IntType)
	serialDef := serializeMapFunc(typ)
	assert.Equal(t, SerialHelperPackage, serialDef.PackageName())

	wc := sai.New()
	f := serialHelperPackage().File("serial.gothic")
	f.Package().ImportResolver()
	f.Writer = wc
	f.Prepare()
	f.Generate()

	expectStrs := []string{
		"github.com/adamcolton/gothic/serial",
		"func Marshalmapstringint(s map[string]int) []byte",
		"func Unmarshalmapstringint(b *[]byte) map[string]int",
	}
	got := wc.String()
	for _, str := range expectStrs {
		assert.Contains(t, got, str)
	}

	// clear serial helper package for other tests
	shp = nil
}
