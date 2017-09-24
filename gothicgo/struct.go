package gothicgo

import (
	"fmt"
	"go/format"
	"sort"
	"strings"
)

// Struct represents a Go struct.
type Struct struct {
	name         string
	file         *File
	fields       map[string]*Field
	fieldOrder   []string
	methods      map[string]*Method
	litMethods   bool
	ReceiverName string
}

// NewStruct adds a Struct to a Package, the file for the struct is automatically
// generated
func (p *Package) NewStruct(name string) *Struct {
	return p.File(name + ".gothic").NewStruct(name)
}

// NewStruct adds a new Struct to an existing file
func (f *File) NewStruct(name string) *Struct {
	s := &Struct{
		name:         name,
		file:         f,
		fields:       make(map[string]*Field),
		methods:      make(map[string]*Method),
		ReceiverName: strings.ToLower(string([]rune(name)[0])),
	}
	f.AddGenerators(s)
	return s
}

// Type returns an object that fulfills the Type interface for this Struct
func (s *Struct) Type() StructType { return &sT{s} }

// Ptr returns an object that fulfills the Type interface for a pointer to this
// Struct
func (s *Struct) Ptr() PointerType { return PointerTo(&sT{s}) }

// AsRet is a helper for returning a pointer to the Struct in a funciton or
// method
func (s *Struct) AsRet() *NameType { return Ret(PointerTo(&sT{s})) }

// AsArg is a helper for passing a pointer to the Struct as an argument to a
// funciton or method
func (s *Struct) AsArg(name string) *NameType { return Arg(name, PointerTo(&sT{s})) }

// AsNmRet is a helper for returning a pointer to the Struct in a funciton or
// method as a named return.
func (s *Struct) AsNmRet(name string) *NameType { return NmRet(name, PointerTo(&sT{s})) }

// File getter. Receiver methods can be added to the file or the Package can be
// accessed through the file and receivers can be added to other files in the
// Package.
func (s *Struct) File() *File { return s.file }

// Name of the struct
func (s *Struct) Name() string { return s.name }

// PackageRef gets the name of the package.
func (s *Struct) PackageRef() PackageRef { return s.file.Package() }

// Field returns a field by name
func (s *Struct) Field(name string) (*Field, bool) {
	f, ok := s.fields[name]
	return f, ok
}

// Prefix fulfills Prefixer
func (s *Struct) Prefix(ref PackageRef) string { return s.file.Imports.Prefix(ref) }

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
		key = typ.Name()
	}
	if f, exists := s.fields[key]; exists {
		return f, fmt.Errorf("Field %s already exists in %s", key, s.name)
	}
	f := &Field{
		nameType: &NameType{
			N: name,
			T: typ,
		},
		Tags: make(map[string]string),
		stct: s,
		//SC:   gothic.NewSC(),
	}
	s.fields[key] = f
	s.fieldOrder = append(s.fieldOrder, key)
	return f, nil
}

// Embed a type as a field
func (s *Struct) Embed(typ Type) (*Field, error) {
	return s.AddField("", typ)
}

func (s *Struct) str() string {
	l := make([]string, len(s.fields)+2)
	l[0] = fmt.Sprintf("type %s struct{", s.name)
	i := 1
	for _, f := range s.fieldOrder {
		l[i] = s.fields[f].String()
		i++
	}
	l[i] = "}"
	return strings.Join(l, "\n")
}

// String returns the struct as Go code
func (s *Struct) String() string {
	b, _ := format.Source([]byte(s.str()))
	return string(b)
}

// Prepare adds all the types to the file import
func (s *Struct) Prepare() error {
	for _, f := range s.fields {
		s.file.AddRefImports(f.Type().PackageRef())
	}
	return nil
}

// Generate adds the Struct to the file
func (s *Struct) Generate() error {
	s.file.AddCode(s.str())
	return nil
}

// PtrMethods returns a bool indication if the methods will be defined on the
// struct literal or the struct pointer.
func (s *Struct) PtrMethods(b bool) { s.litMethods = !b }

// MethodType returns the type used when defining methods on the struct, either
// the struct literal or a pointer to the struct.
func (s *Struct) MethodType() Type {
	if s.litMethods {
		return s.Type()
	}
	return PointerTo(s.Type())
}

// NewMethod on the struct
func (s *Struct) NewMethod(name string, args ...*NameType) *Method {
	m := &Method{
		Ptr:          !s.litMethods,
		ReceiverName: s.ReceiverName,
		strct:        s,
		Func:         NewFunc(name, args...),
	}
	m.Func.File = s.File()
	m.Imports = m.Func.File.Imports
	s.File().AddGenerators(m)
	s.methods[name] = m
	return m
}

// Method gets a method by name
func (s *Struct) Method(name string) (*Method, bool) {
	m, ok := s.methods[name]
	return m, ok
}

// Field is a struct field. Tags follows the convention of `key1:"value1"
// key2:"value2"`. If no value is defined only the key is printed.
type Field struct {
	nameType *NameType
	Tags     map[string]string
	stct     Prefixer
}

// Name of the field. For an embedded field, this will be an empty string.
func (f *Field) Name() string { return f.nameType.Name() }

// Type of the field
func (f *Field) Type() Type { return f.nameType.Type() }

// String returns Go code for the field
func (f *Field) String() string {
	tags := ""
	if len(f.Tags) > 0 {
		s := make([]string, len(f.Tags))
		var i int
		for k, v := range f.Tags {
			if v == "" {
				s[i] = fmt.Sprintf("%s", k)
			} else {
				s[i] = fmt.Sprintf("%s:\"%s\"", k, v)
			}
			i++
		}
		sort.Strings(s)
		tags = " `" + strings.Join(s, " ") + "`"
	}
	typeString := f.Type().RelStr(f.stct)
	if f.Name() == "" {
		return typeString + tags
	}
	return f.Name() + " " + typeString + tags
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

func (s *sT) Name() string           { return s.S.Name() }
func (s *sT) String() string         { return s.S.PackageRef().String() + "." + s.S.Name() }
func (s *sT) File() *File            { return s.S.File() }
func (s *sT) PackageRef() PackageRef { return s.S.PackageRef() }
func (s *sT) Kind() Kind             { return StructKind }
func (s *sT) RelStr(p Prefixer) string {
	return p.Prefix(s.S.PackageRef()) + s.S.Name()
}

type structT struct {
	ref  PackageRef
	name string
}

func (s *structT) Name() string             { return s.name }
func (s *structT) String() string           { return s.ref.Name() + "." + s.name }
func (s *structT) File() *File              { return nil }
func (s *structT) RelStr(p Prefixer) string { return p.Prefix(s.ref) + s.name }
func (s *structT) PackageRef() PackageRef   { return s.ref }
func (s *structT) Kind() Kind               { return StructKind }

// DefStruct returns a StructType for a struct in a package.
func DefStruct(ref PackageRef, name string) StructType {
	return &structT{
		ref:  ref,
		name: name,
	}
}

// Method on a struct
type Method struct {
	*Func
	Ptr          bool
	ReceiverName string
	strct        *Struct
}

// SetName of the method, also updates the method map in the struct.
func (m *Method) SetName(name string) {
	delete(m.strct.methods, m.Func.Name())
	m.Func.SetName(name)
	m.strct.methods[name] = m
}

// String outputs the entire function as a string
func (m *Method) String() string {
	str, _ := m.str()
	return str
}

func (m *Method) str() (string, error) {
	body, err := m.Func.Body()
	if err != nil {
		return "", err
	}
	s := make([]string, 13)
	s[0] = "func ("
	s[1] = m.ReceiverName
	if m.Ptr {
		s[2] = " *"
	} else {
		s[2] = " "
	}
	s[3] = m.strct.Name()
	s[4] = ") "
	s[5] = m.Func.Name()
	s[6] = "("
	s[7] = nameTypeSliceToString(m.strct, m.Func.Args, m.Func.Variadic)
	if l := len(m.Func.Rets); l > 1 || (l == 1 && m.Func.Rets[0].N != "") {
		s[8] = ") ("
		s[10] = ") {\n"
	} else {
		s[8] = ") "
		s[10] = " {\n"
	}
	s[9] = nameTypeSliceToString(m.strct, m.Func.Rets, false)
	s[11] = body
	s[12] = "\n}\n\n"
	return strings.Join(s, ""), nil
}

// Generate writes the method to the file
func (m *Method) Generate() error {
	str, err := m.str()
	if err != nil {
		return err
	}
	m.strct.file.AddCode(str)
	return nil
}
