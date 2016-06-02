package entityBP

import (
	"github.com/adamcolton/gothic/blueprint"
	"github.com/adamcolton/gothic/serial/serialBP"
	"github.com/adamcolton/gothic/structBP"
)

type privStruct struct{ *structBP.Struct } // makes Struct unexported within Entity
type Entity struct {
	privStruct
	ImportHelpers bool
	ImportSerial  bool
	validators    map[string][]Validator
}

func New(pkg, name string) *Entity {
	ent := &Entity{
		privStruct:    privStruct{structBP.New(pkg, name)},
		validators:    map[string][]Validator{},
		ImportHelpers: true,
		ImportSerial:  true,
	}

	serialBP.RegisterSerializeFuncs(ent.String(), serialBP.SerializeFuncs{
		MarshalStr:   "%s.Marshal()",
		UnmarshalStr: "Unmarshal" + name + "(%s)",
		Package:      pkg,
		Marshaler:    serialBP.SimpleMarhsal,
		Unmarshaler:  serialBP.PrependPkgUnmarshal,
	})
	serialBP.RegisterSerializeFuncs(ent.String()+"Ref", serialBP.SerializeFuncs{
		MarshalStr:   "%s.Marshal()",
		UnmarshalStr: "Unmarshal" + name + "Ref(%s)",
		Package:      pkg,
		Marshaler:    serialBP.SimpleMarhsal,
		Unmarshaler:  serialBP.PrependPkgUnmarshal,
	})

	ent.
		ImportGothic("entity").
		Embed(blueprint.TypeString("entity.Ent"))

	ent.AddGenerator(ent)

	return ent
}

func (ent *Entity) Field(name string, t blueprint.Type, validators ...Validator) *Entity {
	ent.privStruct.Field(name, t)
	ent.validators[name] = validators
	return ent
}

func (ent *Entity) Ref() blueprint.Type {
	return blueprint.TypeString(ent.String() + "Ref")
}

// The following are wrapped to allow chaining
func (ent *Entity) Import(imports ...string) *Entity {
	ent.Importer.Import(imports...)
	return ent
}

func (ent *Entity) ImportGothic(imports ...string) *Entity {
	ent.Importer.ImportGothic(imports...)
	return ent
}
func (ent *Entity) ImportApp(imports ...string) *Entity {
	ent.Importer.ImportApp(imports...)
	return ent
}

func (ent *Entity) AddGenerator(gen blueprint.Generator) *Entity {
	ent.BP.AddGenerator(gen)
	return ent
}

func init() {
	serialBP.RegisterSerializeFuncs("entity.Ent", serialBP.SerializeFuncs{
		MarshalStr:   "%s.Marshal()",
		UnmarshalStr: "entity.Unmarshal(%s)",
		Package:      "entity",
		Marshaler:    serialBP.SimpleMarhsal,
		Unmarshaler:  serialBP.SimpleUnmarhsal,
	})
}
