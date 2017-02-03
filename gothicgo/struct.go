package gothicgo

import (
	"fmt"
	"github.com/adamcolton/gothic"
	"go/format"
	"strings"
)

type Struct struct {
	name       string
	file       *File
	fields     map[string]*Field
	fieldOrder []string
	methods    map[string]*Method
	litMethods bool
}

// Adds a Struct to a Package, the file for the struct is automatically
// generated
func (p *Package) NewStruct(name string) *Struct {
	return p.File(name + ".gothic").NewStruct(name)
}

// Adds a new Struct to an existing file
func (f *File) NewStruct(name string) *Struct {
	s := &Struct{
		name:    name,
		file:    f,
		fields:  make(map[string]*Field),
		methods: make(map[string]*Method),
	}
	f.AddFragGen(s)
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

// AsRet is a helper for returning a pointer to the Struct in a funciton or
// method as a named return.
func (s *Struct) AsNmRet(name string) *NameType { return NmRet(name, PointerTo(&sT{s})) }

// Getter for the file. Receiver methods can be added to the file or the Package
// can be accessed through the file and receivers can be added to other files
// in the Package.
func (s *Struct) File() *File         { return s.file }
func (s *Struct) Name() string        { return s.name }
func (s *Struct) PackageName() string { return s.file.Package().Name }

func (s *Struct) Field(name string) (*Field, bool) {
	f, ok := s.fields[name]
	return f, ok
}

func (s *Struct) Fields() []string { return s.fieldOrder }

func (s *Struct) AddField(name string, typ Type) *Field {
	key := name
	if name == "" {
		key = typ.Name()
	}
	f := &Field{
		nameType: &NameType{
			N: name,
			T: typ,
		},
		tags: map[string]string{},
		stct: s.Type(),
		SC:   gothic.NewSC(),
	}
	s.fields[key] = f
	s.fieldOrder = append(s.fieldOrder, key)
	return f
}

func (s *Struct) Embed(typ Type) *Field {
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

func (s *Struct) String() string {
	b, _ := format.Source([]byte(s.str()))
	return string(b)
}

func (s *Struct) Prepare() {
	for _, f := range s.fields {
		s.file.AddPackageImport(f.Type().PackageName())
	}
}

func (s *Struct) Generate() []string {
	return []string{s.str()}
}

func (s *Struct) PtrMethods(b bool) { s.litMethods = !b }

func (s *Struct) NewMethod(name string, args ...*NameType) *Method {
	m := &Method{
		Ptr:          !s.litMethods,
		ReceiverName: strings.ToLower(string([]rune(s.name)[0])), // lower case, first char
		strct:        s,
		Func:         NewFunc(name, args...),
	}
	m.Func.File = s.File()
	s.File().AddFragGen(m)
	s.methods[name] = m
	return m
}

func (s *Struct) Method(name string) (*Method, bool) {
	m, ok := s.methods[name]
	return m, ok
}

type Field struct {
	nameType *NameType
	tags     map[string]string
	stct     interface { // this makes Field easier to test
		PackageName() string
	}
	*gothic.SC
}

func (f *Field) Name() string { return f.nameType.Name() }
func (f *Field) Type() Type   { return f.nameType.Type() }
func (f *Field) Tag(key string) (string, bool) {
	t, ok := f.tags[key]
	return t, ok
}
func (f *Field) SetTag(key, val string) { f.tags[key] = val }
func (f *Field) String() string {
	tags := ""
	if len(f.tags) > 0 {
		s := make([]string, len(f.tags)+2)
		s[0] = " `"
		i := 1
		for k, v := range f.tags {
			s[i] = fmt.Sprintf("%s:\"%s\"", k, v)
			i++
		}
		s[i] = "`"
		tags = strings.Join(s, "")
	}
	typeString := f.Type().RelStr(f.stct.PackageName())
	if f.Name() == "" {
		return typeString + tags
	}
	return f.Name() + " " + typeString + tags
}

type StructType interface {
	Type
}

type sT struct {
	S interface {
		Name() string
		PackageName() string
		File() *File
	}
}

func (s *sT) Name() string        { return s.S.Name() }
func (s *sT) String() string      { return s.S.PackageName() + "." + s.S.Name() }
func (s *sT) File() *File         { return s.S.File() }
func (s *sT) PackageName() string { return s.S.PackageName() }
func (s *sT) Kind() Kind          { return StructKind }
func (s *sT) RelStr(pkg string) string {
	if pkg == s.S.PackageName() {
		return s.S.Name()
	}
	return s.S.PackageName() + "." + s.S.Name()
}

type StructT struct {
	pkgName string
	name    string
}

func (s *StructT) Name() string        { return s.name }
func (s *StructT) String() string      { return s.pkgName + "." + s.name }
func (s *StructT) File() *File         { return nil }
func (s *StructT) PackageName() string { return s.pkgName }
func (s *StructT) Kind() Kind          { return StructKind }
func (s *StructT) RelStr(pkg string) string {
	if pkg == s.pkgName {
		return s.name
	}
	return s.pkgName + "." + s.name
}

func DefStruct(signature string) StructType {
	s := strings.Split(signature, ".")
	if len(s) > 1 {
		return &StructT{
			pkgName: s[0],
			name:    s[1],
		}
	}
	return &StructT{
		name: s[0],
	}
}

type Method struct {
	*Func
	Ptr          bool
	ReceiverName string
	strct        *Struct
}

// String outputs the entire function as a string
func (m *Method) String() string {
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
	s[5] = m.Func.Name
	s[6] = "("
	s[7] = nameTypeSliceToString(m.Func.Args, m.strct.PackageName(), m.Func.Variadic)
	if l := len(m.Func.Rets); l > 1 || (l == 1 && m.Func.Rets[0].N != "") {
		s[8] = ") ("
		s[10] = ") {\n"
	} else {
		s[8] = ")"
		s[10] = " {\n"
	}
	s[9] = nameTypeSliceToString(m.Func.Rets, m.strct.PackageName(), false)
	s[11] = m.Func.Body
	s[12] = "\n}\n\n"
	return strings.Join(s, "")
}

func (m *Method) Generate() []string {
	return []string{m.String()}
}
