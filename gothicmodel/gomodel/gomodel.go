package gomodel

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel"
)

func Struct(pkg *gothicgo.Package, model *gothicmodel.Model) *GoModel {
	s := pkg.NewStruct(model.Name())
	for _, field := range model.Fields() {
		t, ok := Types[field.Type()]
		if !ok {
			continue
		}
		s.AddField(field.Name(), t)
	}
	return &GoModel{
		Struct: s,
		Model:  model,
	}
}

type GoModel struct {
	*gothicgo.Struct
	Model *gothicmodel.Model
}

func (g *GoModel) Fields() []Field {
	gm_fs := g.Model.Fields()
	gg_fs := make([]Field, 0, g.Struct.FieldCount())
	for _, f := range gm_fs {
		kind, ok := Types[f.Type()]
		if !ok {
			continue
		}
		gg_fs = append(gg_fs, Field{
			base: f,
			kind: kind,
		})
	}
	return gg_fs
}

func (g *GoModel) Field(name string) (Field, bool) {
	var field Field
	mf, ok := g.Model.Field(name)
	if !ok {
		return field, false
	}
	gt, ok := Types[mf.Type()]
	if !ok {
		return field, false
	}
	field.base = mf
	field.kind = gt
	return field, true
}

type Field struct {
	base gothicmodel.Field
	kind gothicgo.Type
}

func (f Field) Name() string          { return f.base.Name() }
func (f Field) Type() string          { return f.base.Type() }
func (f Field) Primary() bool         { return f.base.Primary() }
func (f Field) GoType() gothicgo.Type { return f.kind }
