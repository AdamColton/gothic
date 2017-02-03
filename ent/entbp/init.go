package entbp

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicserial"
	_ "github.com/adamcolton/gothic/serial/serialbp"
)

var EntType = gothicgo.DefStruct("ent.Ent")

func init() {
	gothicserial.RegisterSerializeDef("ent.Ent", gothicserial.SerializeDef(&gothicserial.SerializeFuncs{
		MarshalStr:   "%s.Marshal()",
		UnmarshalStr: "ent.Unmarshal(%s)",
		Marshaler:    gothicserial.SimpleMarhsal,
		Unmarshaler:  gothicserial.SimpleUnmarhsal,
		PkgName:      "ent",
	}))
}
