package gomodel

import (
	"github.com/adamcolton/gothic/gothicgo"
	"github.com/adamcolton/gothic/gothicmodel"
)

func Struct(pkg *gothicgo.Package, model *gothicmodel.Model) *GoModel {
	s := pkg.NewStruct(model.Name())
	for _, f := range model.Fields() {
		ts, _ := model.Field(f)
		t, ok := Types[ts]
		if !ok {
			continue
		}
		s.AddField(f, t)
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
