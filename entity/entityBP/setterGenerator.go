package entityBP

import (
	"github.com/adamcolton/gothic/blueprint"
	"github.com/adamcolton/gothic/serial/serialBP"
)

type Setter struct {
	name   string
	on     *Entity
	fields []string
}

func (ent *Entity) Set(name string) *Setter {
	method := &Setter{
		name:   name,
		on:     ent,
		fields: make([]string, 0),
	}
	ent.AddGenerator(method)
	return method
}

func (setter *Setter) Fields(fields ...string) *Setter {
	setter.fields = append(setter.fields, fields...)
	return setter
}

func (setter *Setter) Prepare() {}

const (
	setterHeader     = "func (o *%s) %s(b *[]byte) validation.Errs {\n  valErrs := validation.Errors()"
	setterExtractRow = "  %s := %s"
	setterCheckErrs  = "  if valErrs.HasErrs(){\n    return valErrs\n  }"
	setterSetRow     = "  o.%s = %s"
)

func (setter *Setter) Export() string {
	str := blueprint.StringBuilderLn(setterHeader, setter.on.Name(), setter.name)
	for _, fn := range setter.fields {
		f, _ := setter.on.ByName(fn)
		sf := serialBP.Serialize(f.Type())
		str.AddLn(setterExtractRow, fn, sf.Unmarshal("b", setter.on.Package()))
		validators := setter.on.validators[f.ID()]
		for _, v := range validators {
			str.AddLn(v.Export(fn, setter.on))
		}
	}
	str.AddLn(setterCheckErrs)
	for _, fn := range setter.fields {
		str.AddLn(setterSetRow, fn, fn)
	}
	str.AddLn("  return valErrs\n}")
	return str.String()
}
