package htmlmodel

import (
	"github.com/adamcolton/gothic/gothichtml"
	"github.com/adamcolton/gothic/gothicmodel"
)

type TypeMap map[string]string

var InputTypes = TypeMap{
	"bool":     "checkbox",
	"int":      "number",
	"string":   "text",
	"password": "password",
}

type FieldGenerator func(field, kind string, model *gothicmodel.Model) gothichtml.Node

func (tm TypeMap) GenerateFields(model *gothicmodel.Model, generator FieldGenerator, container gothichtml.ContainerNode) {
	for _, f := range model.Fields() {
		kind, _ := model.Field(f)
		if htmlKind, ok := tm[kind]; ok {
			container.AddChildren(generator(f, htmlKind, model))
		}
	}
}
