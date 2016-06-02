package entityBP

import (
	"github.com/adamcolton/gothic/blueprint"
	"github.com/adamcolton/gothic/serial/serialBP"
)

type marshalGenerator struct {
	ent  *Entity
	rows []string
}

func (ent *Entity) addMarshalGenerator() *Entity {
	return ent.AddGenerator(&marshalGenerator{ent, []string{}})
}

const (
	marshalStructHeader = "func (o *%s) Marshal() []byte{\n  b := []byte{}\n  if o != nil {"
	marshalStructRow    = "    b = append(b, %s...)"
)

func (mg *marshalGenerator) Export() string {
	str := blueprint.StringBuilderLn(marshalStructHeader, mg.ent.Name())
	for _, row := range mg.rows {
		str.AddLn(row)
	}
	str.AddLn("  }\n  return b\n}")
	return str.String()
}

func (mg *marshalGenerator) Prepare() {
	mg.rows = make([]string, len(mg.ent.Fields()))
	for i, f := range mg.ent.Fields() {
		name := "o." + f.ID()
		sf := serialBP.Serialize(f.Type())
		mg.rows[i] = blueprint.StringBuilderLn(marshalStructRow, sf.Marshal(name, mg.ent.Package())).String()
	}
}
