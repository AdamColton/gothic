package gothicgo

import (
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"strings"
)

type ExternalType struct {
	Ref  PackageRef
	Name string
	Type
}

func (e *ExternalType) String() string { return e.Ref.Name() + "." + e.Name }
func (e *ExternalType) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString(p.Prefix(e.Ref))
	sw.WriteString(e.Name)
	sw.Err = errCtx(sw.Err, "While writing external type %s", e.Name)
	return sw.Rets()
}
func (e *ExternalType) PackageRef() PackageRef { return e.Ref }
func (e *ExternalType) Kind() Kind {
	return TypeDefKind
}

// NewExternalType returns a StructType for a struct in a package.
func NewExternalType(ref PackageRef, name string) *ExternalType {
	return &ExternalType{
		Ref:  ref,
		Name: name,
	}
}
func (e *ExternalType) StructEmbedName() string { return e.Name }

func (e *ExternalType) RegisterImports(i *Imports) {
	i.AddRefImports(e.Ref)
}

func (e *ExternalType) Named(name string) NameType {
	return NameType{name, e}
}

type ExternalFunc struct {
	Ref PackageRef
	FuncSig
}

func NewExternalFunc(pkg PackageRef, name string, args ...NameType) *ExternalFunc {
	return &ExternalFunc{
		Ref:     pkg,
		FuncSig: NewFuncSig(name, args...),
	}
}

// Call produces a invocation of the function and fulfills the FuncCaller
// interface
func (e *ExternalFunc) Call(pre Prefixer, args ...string) string {
	buf := bufpool.Get()
	buf.WriteString(pre.Prefix(e.Ref))
	buf.WriteString(e.Name)
	buf.WriteRune('(')
	buf.WriteString(strings.Join(args, ", "))
	buf.WriteRune(')')
	str := buf.String()
	bufpool.Put(buf)
	return str
}

type ExternalInterfaceType struct {
	pkg  PackageRef
	name string
}

func NewExternalInterfaceType(pkg PackageRef, name string) *ExternalInterfaceType {
	return &ExternalInterfaceType{
		pkg:  pkg,
		name: name,
	}
}

func (i *ExternalInterfaceType) String() string { return i.pkg.Name() + "." + i.name }
func (i *ExternalInterfaceType) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString(pre.Prefix(i.pkg))
	sw.WriteString(i.name)
	sw.Err = errCtx(sw.Err, "While writing external interface reference %s", i.name)
	return sw.Rets()
}
func (i *ExternalInterfaceType) PackageRef() PackageRef { return i.pkg }
func (i *ExternalInterfaceType) Kind() Kind             { return InterfaceTypeDefKind }

func (i *ExternalInterfaceType) RegisterImports(im *Imports) {
	im.AddRefImports(i.pkg)
}

func (i *ExternalInterfaceType) InterfaceEmbedName() string {
	return i.name
}
