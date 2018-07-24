package gothicgo

import (
	"fmt"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"sort"
)

// Struct represents a Go struct.
type Struct struct {
	fields     map[string]*Field
	fieldOrder []string
}

// MustStruct adds a new Struct to an existing file and panics if it fails
func NewStruct() *Struct {
	return &Struct{
		fields: make(map[string]*Field),
	}
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

// File getter. Receiver methods can be added to the file or the Package can be
// accessed through the file and receivers can be added to other files in the
// Package.
func (s *Struct) File() *File { return nil }

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
func (s *Struct) AddField(name string, typ Type) (*Field, error) {
	key := name
	if name == "" {
		if emb, ok := typ.(StructEmbeddable); ok {
			key = emb.StructEmbedName()
		} else {
			return nil, fmt.Errorf("Cannot Embed Type: ", typ.String())
		}
	}
	if f, exists := s.fields[key]; exists {
		return f, fmt.Errorf("Field %s already exists in stuct", key)
	}
	f := &Field{
		nameType: NameType{
			N: name,
			T: typ,
		},
		Tags: make(map[string]string),
	}
	s.fields[key] = f
	s.fieldOrder = append(s.fieldOrder, key)
	return f, nil
}

// Embed a type as a field
func (s *Struct) Embed(typ Type) (*Field, error) {
	return s.AddField("", typ)
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
	return sum.Sum, sum.Err
}

// String returns the struct as Go code
func (s *Struct) String() string {
	return typeToString(s, DefaultPrefixer)
}

// Field is a struct field. Tags follows the convention of `key1:"value1"
// key2:"value2"`. If no value is defined only the key is printed.
type Field struct {
	nameType NameType
	Tags     map[string]string
}

// Name of the field. For an embedded field, this will be an empty string.
func (f *Field) Name() string { return f.nameType.Name() }

// Type of the field
func (f *Field) Type() Type { return f.nameType.Type() }

// String returns Go code for the field
func (f *Field) String() string {
	return typeToString(f, DefaultPrefixer)
}

// WriteTo writes the Field to the writer
func (f *Field) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sum := gothicio.NewSumWriter(w)
	if name := f.Name(); name != "" {
		sum.WriteString(name)
		sum.WriteString(" ")
	}
	f.Type().PrefixWriteTo(sum, p)

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

	sum.Err = errCtx(sum.Err, "While writing field %s", f.nameType.N)

	return sum.Sum, sum.Err
}

// StructType is just a wrapper around Type
type StructType interface {
	Type
}

type sT struct {
	S interface {
		Name() string
		PackageRef() PackageRef
		File() *File
	}
}

func (s *sT) String() string         { return s.S.PackageRef().String() + "." + s.S.Name() }
func (s *sT) File() *File            { return s.S.File() }
func (s *sT) PackageRef() PackageRef { return s.S.PackageRef() }
func (s *sT) Kind() Kind             { return StructKind }
func (s *sT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString(p.Prefix(s.S.PackageRef()))
	sw.WriteString(s.S.Name())
	return sw.Rets()
}

type structT struct {
	ref  PackageRef
	name string
}

func (s *structT) String() string { return s.ref.Name() + "." + s.name }
func (s *structT) File() *File    { return nil }
func (s *structT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString(p.Prefix(s.ref))
	sw.WriteString(s.name)
	return sw.Rets()
}
func (s *structT) PackageRef() PackageRef { return s.ref }
func (s *structT) Kind() Kind             { return StructKind }

// DefStruct returns a StructType for a struct in a package.
func DefStruct(ref PackageRef, name string) StructType {
	return &structT{
		ref:  ref,
		name: name,
	}
}
func (s *structT) StructEmbedName() string { return s.name }

func (s *structT) RegisterImports(i *Imports) {
	i.AddRefImports(s.ref)
}
