package entbp

import (
	"github.com/adamcolton/gothic/gothicgo"
)

// EntBP represents an entity.
type EntBP struct {
	name       string
	fields     map[string]*Field
	fieldOrder []string
	pkg        Package
	ref        gothicgo.StructType
}

// AddField adds a field to an entity defined by name and the type in Go. A
// reference to the EntBP is returned so that AddField can be chained.
func (e *EntBP) AddField(name string, typ gothicgo.Type) *EntBP {
	key := name
	if name == "" {
		key = typ.Name()
	}
	f := &Field{
		name: name,
		typ:  typ,
		id:   len(e.fields),
	}
	e.fields[key] = f
	e.fieldOrder = append(e.fieldOrder, key)
	return e
}

// Field represents an entity field. It has a name, (Go) type and ID. The ID
// is used when renaming fields.
type Field struct {
	name string
	typ  gothicgo.Type
	id   int
}

// Ref returns a Go reference type for this entity. A reference type is like a
// pointer for entities - it is a reference to another entity which is
// marshalled by id rather than value.
func (e *EntBP) Ref() gothicgo.StructType {
	if e.ref == nil {
		e.ref = gothicgo.PointerTo(gothicgo.DefStruct(e.pkg.Name() + "." + e.name + "Ref"))
	}
	return e.ref
}

// Package returns the package to which the entity belongs.
func (e *EntBP) Package() Package { return e.pkg }
