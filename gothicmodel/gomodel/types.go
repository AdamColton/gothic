package gomodel

import (
	"github.com/adamcolton/gothic/gothicgo"
)

const (
	Bool       = "bool"
	Byte       = "byte"
	Int        = "int"
	Int8       = "int8"
	Int16      = "int16"
	Int32      = "int32"
	Int64      = "int64"
	Complex128 = "complex128"
	Complex64  = "complex64"
	Float32    = "float32"
	Float64    = "float64"
	Rune       = "rune"
	String     = "string"
	Uint       = "uint"
	Uint8      = "uint8"
	Uint16     = "uint16"
	Uint32     = "uint32"
	Uint64     = "uint64"
	Uintptr    = "uintptr"
	Datetime   = "datetime"
)

// Types maps model Type stirngs to go types. This should be extended as new
// types are defined.
var Types = map[string]gothicgo.Type{
	Bool:       gothicgo.BoolType,
	Byte:       gothicgo.ByteType,
	Int:        gothicgo.IntType,
	Int8:       gothicgo.Int8Type,
	Int16:      gothicgo.Int16Type,
	Int32:      gothicgo.Int32Type,
	Int64:      gothicgo.Int64Type,
	Complex128: gothicgo.Complex128Type,
	Complex64:  gothicgo.Complex64Type,
	Float32:    gothicgo.Float32Type,
	Float64:    gothicgo.Float64Type,
	Rune:       gothicgo.RuneType,
	String:     gothicgo.StringType,
	Uint:       gothicgo.UintType,
	Uint8:      gothicgo.Uint8Type,
	Uint16:     gothicgo.Uint16Type,
	Uint32:     gothicgo.Uint32Type,
	Uint64:     gothicgo.Uint64Type,
	Uintptr:    gothicgo.UintptrType,
	Datetime:   gothicgo.DefStruct(gothicgo.MustPackageRef("time"), "Time"),
}

// Tags defines the tags that will set on a field
type Tags map[string]string

// TypeTags defines the tags that will set on fields
var TypeTags = make(map[string]Tags)
