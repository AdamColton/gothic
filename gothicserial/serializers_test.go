package gothicserial

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterSerializeDef(t *testing.T) {
	p := gothicgo.NewPackage("test")
	sd := SerializeFuncs{
		MarshalStr:   "test",
		UnmarshalStr: "test",
		PkgName:      "test",
		F:            p.File("test"),
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	}
	RegisterSerializeDef("test", sd.Def())
	d, ok := GetSerializeDef("test")
	assert.True(t, ok)
	assert.Equal(t, sd.Marshal("a", "pkg"), d.Marshal("a", "pkg"))
}
