package entityBP

import (
	"github.com/adamcolton/gothic/blueprint"
	"github.com/adamcolton/gothic/serial/serialBP"
)

type unmarshalGenerator struct {
	ent  *Entity
	rows []string
}

func (ent *Entity) addUnmarshalGenerator() *Entity {
	return ent.AddGenerator(&unmarshalGenerator{
		ent:  ent,
		rows: []string{},
	})
}

const (
	unmarshalStructHeader = "func Unmarshal%s(b *[]byte) *%s {\n  return &%s{"
	unmarshalStructRow    = "%s: %s,"
)

func (ug *unmarshalGenerator) Export() string {
	str := blueprint.StringBuilderLn(unmarshalStructHeader, ug.ent.Name(), ug.ent.Name(), ug.ent.Name())
	for _, row := range ug.rows {
		str.AddLn(row)
	}
	str.AddLn("  }\n}")
	return str.String()
}

func (ug *unmarshalGenerator) Prepare() {
	ug.rows = make([]string, len(ug.ent.Fields()))
	for i, f := range ug.ent.Fields() {
		sf := serialBP.Serialize(f.Type())
		ug.rows[i] = blueprint.StringBuilderLn(unmarshalStructRow, f.ID(), sf.Unmarshal("b", ug.ent.Package())).String()
	}
}
