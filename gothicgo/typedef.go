package gothicgo

import (
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
	methods      map[string]*TypeDefMethod
	Ptr          bool
}

func (f *File) NewTypeDef(name string, t Type) *TypeDef {
	td := &TypeDef{
		baseType:     t,
		name:         name,
		file:         f,
		methods:      make(map[string]*TypeDefMethod),
		ReceiverName: strings.ToLower(string([]rune(name)[0])),
		Ptr:          true,
	}
	f.AddWriterTo(td)
	return td
}

func (td *TypeDef) WriteTo(w io.Writer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
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

func (td *TypeDef) StructEmbedName() string {
	return td.name
}

// NewMethod on the struct
func (td *TypeDef) NewMethod(name string, args ...NameType) *TypeDefMethod {
	m := &TypeDefMethod{
		typeDef: td,
		Ptr:     td.Ptr,
		Func:    NewFunc(td.File().Imports, name, args...),
	}
	m.Func.File = td.File()
	td.File().AddWriterTo(m)
	td.methods[name] = m
	return m
}

// Method gets a method by name
func (td *TypeDef) Method(name string) (*TypeDefMethod, bool) {
	m, ok := td.methods[name]
	return m, ok
}

// TypeDefMethod on a struct
type TypeDefMethod struct {
	*Func
	Ptr     bool
	typeDef *TypeDef
}

// SetName of the method, also updates the method map in the struct.
func (m *TypeDefMethod) SetName(name string) {
	delete(m.typeDef.methods, m.Func.Name())
	m.Func.Sig.Name = name
	m.typeDef.methods[name] = m
}

// String outputs the entire function as a string
func (m *TypeDefMethod) String() string {
	buf := bufpool.Get()
	m.WriteTo(buf)
	return bufpool.PutStr(buf)
}

// WriteTo writes the TypeDefMethod to the writer
func (m *TypeDefMethod) WriteTo(w io.Writer) (int64, error) {
	sum := gothicio.NewSumWriter(w)
	if m.Comment != "" {
		NewComment(strings.Join([]string{m.Name(), m.Comment}, " ")).WriteTo(sum)
	}
	sum.WriteString("func (")
	m.Receiver().PrefixWriteTo(sum, m.typeDef.file)
	sum.WriteString(") ")
	sum.WriteString(m.Func.Name())
	writeArgsRets(sum, m.typeDef.file, m.Func.Sig.Args, m.Func.Sig.Rets, m.Func.Variadic)
	sum.WriteString("{\n\t")
	if m.Func.Body != nil {
		m.Func.Body.WriteTo(sum)
	}
	sum.WriteString("\n}")

	sum.Err = errCtx(sum.Err, "While writing method %s:", m.Sig.Name)

	return sum.Sum, sum.Err
}

func (m *TypeDefMethod) Receiver() NameType {
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
