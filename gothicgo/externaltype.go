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
