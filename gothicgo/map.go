package gothicgo

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

// MapType extends the Type interface with Map specific information
type MapType interface {
	Type
	Elem() Type
	Key() Type
}

type mapT struct {
	key  Type
	elem Type
}

func (m *mapT) String() string {
	return typeToString(m, DefaultPrefixer)
}

func (m *mapT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString("map[")
	m.key.PrefixWriteTo(sw, p)
	sw.WriteRune(']')
	m.elem.PrefixWriteTo(sw, p)
	return sw.Rets()
}

// PackageRef will always return PkgBuiltin() on Map. Packages of the key and element can
// be inspected independantly
func (m *mapT) PackageRef() PackageRef { return pkgBuiltin }
func (m *mapT) File() *File            { return nil }
func (m *mapT) Kind() Kind             { return MapKind }
func (m *mapT) Elem() Type             { return m.elem }
func (m *mapT) Key() Type              { return m.key }

func (m *mapT) RegisterImports(i *Imports) {
	m.elem.RegisterImports(i)
	m.key.RegisterImports(i)
}

// MapOf returns a MapType
func MapOf(key, elem Type) MapType {
	return &mapT{
		key:  key,
		elem: elem,
	}
}
