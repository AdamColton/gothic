package gothicserial

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterSerializeDef(t *testing.T) {
	ctx := New()
	p := gothicgo.NewPackage("test")
	sd := SerializeFuncs{
		MarshalStr:   "test",
		UnmarshalStr: "test",
		PkgName:      "test",
		F:            p.File("test"),
		Marshaler:    SimpleMarhsal,
		Unmarshaler:  SimpleUnmarhsal,
	}
	ctx.Register(gothicgo.IntType, sd.Def())
	d := ctx.Get(gothicgo.IntType)
	assert.Equal(t, sd.Marshal("a", "pkg"), d.Marshal("a", "pkg"))
}
