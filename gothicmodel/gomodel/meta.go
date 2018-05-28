package gomodel

import (
	"github.com/adamcolton/gothic/gothicmodel"
)

type ModelMetaHandler func(name, value string, model *GoModel)

var ModelMetaHandlers = make(map[string]ModelMetaHandler)

type FieldMetaHandler func(name, value string, field gothicmodel.Field, model *GoModel)

var FieldMetaHandlers = make(map[string]FieldMetaHandler)

func MetaToTag(name, value string, field gothicmodel.Field, model *GoModel) {
	sf, _ := model.Struct.Field(field.Name())
	sf.Tags[name] = value
}
