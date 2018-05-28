package gomodel

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel"
)

func Must(pkg *gothicgo.Package, model *gothicmodel.GothicModel) *GoModel {
	gm, err := New(pkg, model)
	if err != nil {
		panic(err)
	}
	return gm
}

// New creates a GoModel from a Gothic model.
func New(pkg *gothicgo.Package, model *gothicmodel.GothicModel) (*GoModel, error) {
	s, err := pkg.NewStruct(model.Name())
	if err != nil {
		return nil, err
	}
	for _, field := range model.Fields() {
		t, ok := Types[field.Type()]
		if !ok {
			continue
		}
		goField, err := s.AddField(field.Name(), t)
		if err != nil {
			return nil, err
		}
		if tags, ok := TypeTags[field.Type()]; ok {
			goField.Tags = make(map[string]string)
			for k, v := range tags {
				goField.Tags[k] = v
			}
		}
	}

	gm := &GoModel{
		Struct:      s,
		GothicModel: model,
	}

	for _, k := range model.Metas() {
		if h, ok := ModelMetaHandlers[k]; ok {
			v, _ := model.Meta(k)
			h(k, v, gm)
		}
	}

	for _, field := range model.Fields() {
		if _, ok := Types[field.Type()]; !ok {
			continue
		}
		for _, k := range field.Metas() {
			if h, ok := FieldMetaHandlers[k]; ok {
				v, _ := field.Meta(k)
				h(k, v, field, gm)
			}
		}
	}
	return gm, nil
}

// GoModel embeds a Struct and includes a reference to the Model that generated
// it
type GoModel struct {
	*gothicgo.Struct
	GothicModel *gothicmodel.GothicModel
}

// Fields lists the Fields on the Model
func (g *GoModel) Fields() []Field {
	gmFs := g.GothicModel.Fields()
	ggFs := make([]Field, 0, g.Struct.FieldCount())
	for _, f := range gmFs {
		kind, ok := Types[f.Type()]
		if !ok {
			continue
		}
		ggFs = append(ggFs, Field{
			base: f,
			kind: kind,
		})
	}
	return ggFs
}

// Field returns a Field by name
func (g *GoModel) Field(name string) (Field, bool) {
	var field Field
	mf, ok := g.GothicModel.Field(name)
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

// Field holds the information for both the Go Struct field an the Model field.
type Field struct {
	base gothicmodel.Field
	kind gothicgo.Type
}

// Name of the field, which is the same in the Model and the Struct
func (f Field) Name() string { return f.base.Name() }

// Type from the model, returns a string.
func (f Field) Type() string { return f.base.Type() }

// Primary returns true if this is the primary field for the model
func (f Field) Primary() bool { return f.base.Primary() }

// GoType returns the type of the field in Go.
func (f Field) GoType() gothicgo.Type { return f.kind }
