package gothicgo

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"sort"
	"strings"
)

// StructEmbeddable is used to embed a named type in a struct. The returned
// string is what the field name will be. So when embedding *foo.Bar, the
// StructEmbedName will be Bar.
type StructEmbeddable interface {
	Type
	StructEmbedName() string
}

// Struct represents a Go struct.
type Struct struct {
	fields     map[string]*Field
	fieldOrder []string
}

// MustStruct adds a new Struct to an existing file and panics if it fails
func NewStruct(fields ...PrefixWriterTo) *Struct {
	s := &Struct{
		fields: make(map[string]*Field),
	}
	s.AddFields(fields...)
	return s
}

// Ptr returns an object that fulfills the Type interface for a pointer to this
// Struct
func (s *Struct) Ptr() PointerType { return PointerTo(s) }

// AsRet is a helper for returning a the Struct in a funciton or
// method
func (s *Struct) AsRet() NameType { return Ret(s) }

func (s *Struct) Kind() Kind { return StructKind }

// AsArg is a helper for passing a pointer to the Struct as an argument to a
// funciton or method
func (s *Struct) AsArg(name string) NameType { return Arg(name, s) }

// AsNmRet is a helper for returning a pointer to the Struct in a funciton or
// method as a named return.
func (s *Struct) AsNmRet(name string) NameType { return NmRet(name, s) }

// PackageRef gets the name of the package.
func (s *Struct) PackageRef() PackageRef { return nil }

// Field returns a field by name
func (s *Struct) Field(name string) (*Field, bool) {
	f, ok := s.fields[name]
	return f, ok
}

func (s *Struct) RegisterImports(i *Imports) {
	for _, f := range s.fields {
		f.Type().RegisterImports(i)
	}
}

// Fields returns the fields in order.
func (s *Struct) Fields() []string {
	fs := make([]string, len(s.fieldOrder))
	copy(fs, s.fieldOrder)
	return fs
}

// FieldCount returns how many fields the struct has
func (s *Struct) FieldCount() int {
	return len(s.fieldOrder)
}

// AddField to the struct
func (s *Struct) AddFields(fields ...PrefixWriterTo) error {
	for _, p := range fields {
		var f *Field
		switch t := p.(type) {
		case *Field:
			f = t
		case NameType:
			f = &Field{NameType: t}
		default:
			if emb, ok := p.(StructEmbeddable); ok {
				f = NewField("", emb)
			} else if emb, ok := p.(InterfaceEmbeddable); ok {
				f = NewField("", emb)
			} else {
				return fmt.Errorf("Given type cannot be converted to struct field")
			}
		}

		key := f.Name()
		if key == "" {
			return fmt.Errorf("Field must either be named or type must be StructEmbeddable or InterfaceEmbeddable")
		}
		if _, exists := s.fields[key]; exists {
			return fmt.Errorf("Field %s already exists in struct", key)
		}
		s.fields[key] = f
		s.fieldOrder = append(s.fieldOrder, key)
	}
	return nil
}

// Embed a type as a field
func (s *Struct) Embed(typ StructEmbeddable) (*Field, error) {
	f := NewField("", typ)
	return f, s.AddFields(f)
}

// WriteTo writes the Struct to the writer
func (s *Struct) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sum := gothicio.NewSumWriter(w)
	sum.WriteString("struct {")
	for _, f := range s.fieldOrder {
		sum.WriteString("\n\t")
		s.fields[f].PrefixWriteTo(sum, p)
	}
	if len(s.fieldOrder) > 0 {
		sum.WriteString("\n}")
	} else {
		sum.WriteString("}")
	}
	sum.Err = errCtx(sum.Err, "While writing struct")
	return sum.Rets()
}

// String returns the struct as Go code
func (s *Struct) String() string {
	return typeToString(s, DefaultPrefixer)
}

// Field is a struct field. Tags follows the convention of `key1:"value1"
// key2:"value2"`. If no value is defined only the key is printed.
type Field struct {
	NameType
	Tags map[string]string
}

func NewField(name string, typ Type, tags ...string) *Field {
	var tagMap map[string]string
	if len(tags) > 0 {
		if len(tags)%2 == 1 {
			tags = append(tags, "")
		}
		tagMap = make(map[string]string, len(tags)/2)
		for i := 0; i < len(tags); i += 2 {
			tagMap[tags[i]] = tags[i+1]
		}
	}
	return &Field{
		NameType: NameType{name, typ},
		Tags:     tagMap,
	}
}

func (f *Field) Name() string {
	if f.N != "" {
		return f.N
	}
	if emb, ok := f.T.(StructEmbeddable); ok {
		return emb.StructEmbedName()
	}
	if emb, ok := f.T.(InterfaceEmbeddable); ok {
		return emb.InterfaceEmbedName()
	}
	return ""
}

// String returns Go code for the field
func (f *Field) String() string {
	return typeToString(f, DefaultPrefixer)
}

func (f *Field) AddTag(key, value string) {
	if f.Tags == nil {
		f.Tags = map[string]string{
			key: value,
		}
		return
	}
	if s, ok := f.Tags[key]; ok {
		f.Tags[key] = s + ";" + value
		return
	}
	f.Tags[key] = value
}

// WriteTo writes the Field to the writer
func (f *Field) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sum := gothicio.NewSumWriter(w)
	if f.N != "" {
		sum.WriteString(f.N)
		sum.WriteString(" ")
	}
	indentWriter := gothicio.ReplacerWriter{
		Writer:   sum,
		Replacer: strings.NewReplacer("\n", "\n\t"),
	}
	f.Type().PrefixWriteTo(indentWriter, p)

	if len(f.Tags) > 0 {
		sum.WriteString(" `")
		tags := make([]string, 0, len(f.Tags))
		for k := range f.Tags {
			tags = append(tags, k)
		}
		sort.Strings(tags)
		for i, tag := range tags {
			if i > 0 {
				sum.WriteString(" ")
			}
			sum.WriteString(tag)
			if v := f.Tags[tag]; v != "" {
				sum.WriteString(":\"")
				sum.WriteString(v)
				sum.WriteString("\"")
			}
		}
		sum.WriteString("`")
	}

	sum.Err = errCtx(sum.Err, "While writing field %s", f.Name())

	return sum.Rets()
}
