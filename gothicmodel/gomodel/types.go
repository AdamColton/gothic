package gomodel

import (
	"github.com/adamcolton/gothic/gothicgo"
)

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
