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

func (tm TypeMap) GenerateFields(model *gothicmodel.Model, generator FieldGenerator, container gothichtml.ContainerNode, fields ...string) {
	if len(fields) == 0 {
		fields = model.Fields()
	}
	for _, f := range fields {
		kind, _ := model.Field(f)
		if htmlKind, ok := tm[kind]; ok {
			container.AddChildren(generator(f, htmlKind, model))
		}
	}
}
