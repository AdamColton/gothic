package gomodel

import (
	"github.com/adamcolton/gothic/gothicgo"
)

// Types maps model Type stirngs to go types. This should be extended as new
// types are defined.
var Types = map[string]gothicgo.Type{
	"bool":       gothicgo.BoolType,
	"byte":       gothicgo.ByteType,
	"int":        gothicgo.IntType,
	"int8":       gothicgo.Int8Type,
	"int16":      gothicgo.Int16Type,
	"int32":      gothicgo.Int32Type,
	"int64":      gothicgo.Int64Type,
	"complex128": gothicgo.Complex128Type,
	"complex64":  gothicgo.Complex64Type,
	"float32":    gothicgo.Float32Type,
	"float64":    gothicgo.Float64Type,
	"rune":       gothicgo.RuneType,
	"string":     gothicgo.StringType,
	"uint":       gothicgo.UintType,
	"uint8":      gothicgo.Uint8Type,
	"uint16":     gothicgo.Uint16Type,
	"uint32":     gothicgo.Uint32Type,
	"uint64":     gothicgo.Uint64Type,
	"uintptr":    gothicgo.UintptrType,
	"datetime":   gothicgo.DefStruct(gothicgo.MustPackageRef("time"), "Time"),
}

// Tags defines the tags that will set on a field
type Tags map[string]string

// TypeTags defines the tags that will set on fields
var TypeTags = make(map[string]Tags)
