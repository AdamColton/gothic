// StructBP is a blueprint for structs. The inheiritance is a little odd.
// A StructBP has an embedded blueprint.BP, that is registered as implementing
// the Blueprint interface. The Struct struct has Prepare and Export() string
// implementing the generator interface. So the Struct registers itself with
// the BP as a generator.

package structBP

import (
	"github.com/adamcolton/gothic/blueprint"
)

const (
	structHeader = "type %s struct{"
	structRow    = "%s %s"
)

type privBP struct{ *blueprint.BP } // so that BP is unexported within Struct
type Struct struct {
	privBP
	fields   []*Field
	fieldMap map[string]*Field
}

func NewUnreg(pkg, name string) *Struct {
	strct := &Struct{
		privBP:   privBP{blueprint.New(pkg, name)},
		fields:   []*Field{},
		fieldMap: map[string]*Field{},
	}
	strct.AddGenerator(strct)
	return strct
}

func New(pkg, name string) *Struct {
	strct := NewUnreg(pkg, name)
	blueprint.Register(strct.BP)
	return strct
}

func (s *Struct) Prepare() {}

func (s *Struct) Export() string {
	sb := blueprint.StringBuilder(structHeader, s.Name())
	for _, f := range s.fields {
		if f.name == "" {
			sb.AddLn(f.typ.String())
		} else {
			sb.AddLn(structRow, f.name, f.typ.String())
		}
	}
	sb.AddLn("}")
	return sb.String()
}

func (s *Struct) Field(name string, t blueprint.Type) *Struct {
	f := &Field{
		name: name,
		typ:  t,
	}
	s.fields = append(s.fields, f)
	s.fieldMap[name] = f
	return s
}

func (s *Struct) Embed(t blueprint.Type) *Struct {
	f := &Field{
		name: "",
		typ:  t,
	}
	s.fields = append([]*Field{f}, s.fields...)
	s.fieldMap[f.typ.Name()] = f
	return s
}

func (s *Struct) ByName(fieldName string) (*Field, bool) {
	f, ok := s.fieldMap[fieldName]
	return f, ok
}

func (s *Struct) Fields() []*Field { return s.fields }

// The following are wrapped to allow chaining

func (s *Struct) Import(imports ...string) *Struct {
	s.Importer.Import(imports...)
	return s
}

func (s *Struct) ImportGothic(imports ...string) *Struct {
	s.Importer.ImportGothic(imports...)
	return s
}
func (s *Struct) ImportApp(imports ...string) *Struct {
	s.Importer.ImportApp(imports...)
	return s
}

func (s *Struct) AddGenerator(gen blueprint.Generator) *Struct {
	s.BP.AddGenerator(gen)
	return s
}
