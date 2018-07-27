package gothicgo

import (
	"fmt"
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"strings"
)

type TypeDef struct {
	baseType     Type
	name         string
	file         *File
	ReceiverName string
	methods      map[string]*Method
	Ptr          bool
	Comment      string
}

type NewTypeDefiner interface {
	NewTypeDef(name string, t Type) *TypeDef
}

// NewStruct adds a Struct to a Package, the file for the struct is automatically
// generated
func (p *Package) NewTypeDef(name string, t Type) (*TypeDef, error) {
	if _, ok := p.names[name]; ok {
		return nil, fmt.Errorf("Cannot define type %s in package %s; name already exists in scope", name, p.Name())
	}
	return p.File(name+".gothic").NewTypeDef(name, t)
}

func (f *File) NewTypeDef(name string, t Type) (*TypeDef, error) {
	if _, ok := f.pkg.names[name]; ok {
		return nil, fmt.Errorf("Cannot define type %s in package %s; name already exists in scope", name, f.pkg.Name())
	}
	td := &TypeDef{
		baseType:     t,
		name:         name,
		file:         f,
		methods:      make(map[string]*Method),
		ReceiverName: strings.ToLower(string([]rune(name)[0])),
		Ptr:          true,
	}
	return td, f.AddWriterTo(td)
}

func (td *TypeDef) Prepare() error {
	td.baseType.RegisterImports(td.File().Imports)
	return nil
}

func (td *TypeDef) WriteTo(w io.Writer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	if td.Comment != "" {
		NewComment(strings.Join([]string{td.name, td.Comment}, " ")).WriteTo(sw)
	}
	sw.WriteString("type ")
	sw.WriteString(td.name)
	sw.WriteRune(' ')
	td.baseType.PrefixWriteTo(sw, td.file)
	return sw.Rets()
}

func (td *TypeDef) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString(p.Prefix(td.file.Package()))
	sw.WriteString(td.name)
	return sw.Rets()
}
func (td *TypeDef) String() string { return typeToString(td, DefaultPrefixer) }
func (td *TypeDef) PackageRef() PackageRef {
	return td.file.Package()
}
func (td *TypeDef) File() *File {
	return td.file
}
func (td *TypeDef) Kind() Kind {
	return TypeDefKind
}
func (td *TypeDef) Name() string {
	return td.name
}

func (td *TypeDef) RegisterImports(i *Imports) {
	i.AddRefImports(td.file.Package())
}

func (td *TypeDef) StructEmbedName() string {
	return td.name
}

// NewMethod on the struct
func (td *TypeDef) NewMethod(name string, args ...NameType) (*Method, error) {
	if name == "" {
		return nil, errStr("Cannot have unnamed method")
	}
	fn := &Func{
		FuncSig: NewFuncSig(name, args...),
		file:    td.file,
	}
	m := &Method{
		typeDef: td,
		Ptr:     td.Ptr,
		Func:    fn,
	}
	err := td.File().AddWriterTo(m)
	if err != nil {
		return nil, err
	}
	td.methods[name] = m
	return m, nil
}

// Method gets a method by name
func (td *TypeDef) Method(name string) (*Method, bool) {
	m, ok := td.methods[name]
	return m, ok
}

// Method on a struct
type Method struct {
	*Func
	Ptr     bool
	typeDef *TypeDef
}

// SetName of the method, also updates the method map in the struct.
func (m *Method) SetName(name string) {
	delete(m.typeDef.methods, m.Func.Name)
	m.Name = name
	m.typeDef.methods[name] = m
}

// String outputs the entire function as a string
func (m *Method) String() string {
	buf := bufpool.Get()
	m.WriteTo(buf)
	return bufpool.PutStr(buf)
}

// WriteTo writes the Method to the writer
func (m *Method) WriteTo(w io.Writer) (int64, error) {
	if m.Name == "" {
		return 0, errStr("Cannot have unnamed method")
	}
	sum := gothicio.NewSumWriter(w)
	if m.Comment != "" {
		NewComment(strings.Join([]string{m.Name, m.Comment}, " ")).WriteTo(sum)
	}

	sum.WriteString("func (")
	sum.WriteString(m.typeDef.ReceiverName)
	if m.Ptr {
		sum.WriteString(" *")
	} else {
		sum.WriteRune(' ')
	}
	sum.WriteString(m.typeDef.Name())
	sum.WriteString(") ")
	sum.WriteString(m.Name)
	sum.WriteRune('(')
	var str string
	str, sum.Err = nameTypeSliceToString(m.typeDef.file, m.Args, m.Variadic)
	sum.WriteString(str)
	end := " {\n\t"
	if len(m.Rets) > 1 {
		sum.WriteString(") (")
		end = ") {\n\t"
	} else {
		sum.WriteString(") ")
	}
	str, sum.Err = nameTypeSliceToString(m.typeDef.file, m.Rets, false)
	sum.WriteString(str)
	sum.WriteString(end)

	if m.Func.Body != nil {
		m.Func.Body.PrefixWriteTo(sum, m.typeDef.file)
	}
	sum.WriteString("\n}")
	sum.Err = errCtx(sum.Err, "While writing method %s:", m.Name)

	return sum.Rets()
}

func (m *Method) Receiver() NameType {
	n := NameType{
		N: m.typeDef.ReceiverName,
	}
	if m.Ptr {
		n.T = PointerTo(m.typeDef)
	} else {
		n.T = m.typeDef
	}
	return n
}
