package entbp

import (
	"github.com/adamcolton/gothic/gothicgo"
)

type EntBP struct {
	name       string
	fields     map[string]*Field
	fieldOrder []string
	pkg        Package
	ref        gothicgo.StructType
}

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

type Field struct {
	name string
	typ  gothicgo.Type
	id   int
}

func (e *EntBP) Ref() gothicgo.StructType {
	if e.ref == nil {
		e.ref = gothicgo.PointerTo(gothicgo.DefStruct(e.pkg.Name() + "." + e.name + "Ref"))
	}
	return e.ref
}

func (e *EntBP) Package() Package { return e.pkg }
