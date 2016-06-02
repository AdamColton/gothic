package entityBP

import (
	"github.com/adamcolton/gothic/blueprint"
	"github.com/adamcolton/gothic/serial/serialBP"
)

func (ent *Entity) Prepare() {
	ent.
		addMarshalGenerator().
		addUnmarshalGenerator().
		addUnmarshalEntityGenerator()

	if ent.ImportHelpers {
		ent.ImportApp(serialBP.SerialHelperPackage)
	}

	if ent.ImportSerial {
		ent.ImportGothic("serial")
	}

	for _, f := range ent.Fields() {
		validators := ent.validators[f.ID()]
		for _, v := range validators {
			v.Import(f.Name(), ent)
		}
	}
}

func (ent *Entity) Export() string {
	return ent.exportGetter() + "\n" + ent.exportRef()
}

const (
	structHeader = "type %s struct{"
	structRow    = "  %s %s"
)

func (ent *Entity) exportStruct() string {
	str := blueprint.StringBuilderLn(structHeader, ent.Name())
	for _, f := range ent.Fields() {
		if f.Name() == "" {
			str.AddLn(f.Type().String())
		} else {
			str.AddLn(structRow, f.Name(), f.Type().String())
		}
	}
	str.AddLn("}")
	return str.String()
}

const getterTemplate = `
func Get%s(id uint64) *%s {
  var ent entity.Entity
  entity.Get(id, &ent, Unmarshal%sEntity)
  return ent.(*%s)
}
`

func (ent *Entity) exportGetter() string {
	return blueprint.StringBuilder(getterTemplate, ent.Name(), ent.Name(), ent.Name(), ent.Name()).String()
}

const refTemplate = `
type %sRef struct{ entity.Ent }

func (r %sRef) Get() *%s {
  return Get%s(r.ID())
}

func (o *%s) Ref() %sRef {
  return %sRef{entity.Def(o.ID())}
}

func Unmarshal%sRef(b *[]byte) %sRef {
  return %sRef{entity.Unmarshal(b)}
}`

func (ent *Entity) exportRef() string {
	return blueprint.StringBuilder(refTemplate, ent.Name(), ent.Name(), ent.Name(), ent.Name(), ent.Name(), ent.Name(), ent.Name(), ent.Name(), ent.Name(), ent.Name()).String()
}
