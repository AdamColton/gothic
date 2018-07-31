package gothicgo

import (
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"strings"
)

// Interface is used to generate an interface
type Interface struct {
	methods []*interfaceMethod
}

// NewInterface adds a new interface to an existing file
func NewInterface() *Interface {
	return &Interface{}
}

// AddMethod to the interface
func (i *Interface) AddMethod(funcSig FuncSig) {
	i.methods = append(i.methods, &interfaceMethod{
		funcSig: funcSig,
		ifc:     i,
	})
}

func typeSliceToString(nts []NameType, pre Prefixer, variadic bool) string {
	l := len(nts)
	var s = make([]string, l)
	l--
	buf := bufpool.Get()
	for i, nt := range nts {
		nt.T.PrefixWriteTo(buf, pre)
		if i == l && variadic {
			s[i] = " ..." + buf.String()
		} else {
			s[i] = buf.String()
		}
		buf.Reset()
	}
	bufpool.Put(buf)
	return strings.Join(s, ", ")
}

// String returns the interface package and name and fulfills the Type interface
func (i *Interface) String() string {
	return typeToString(i, DefaultPrefixer)
}

// WriteTo writes the interface name and package if necessary.
func (i *Interface) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	s := gothicio.NewSumWriter(w)
	s.WriteString("interface {")
	for _, m := range i.methods {
		m.PrefixWriteTo(s, pre)
		s.WriteString("\n\t")
	}
	if len(i.methods) > 0 {
		s.WriteString("\n}")
	} else {
		s.WriteString("}")
	}
	s.Err = errCtx(s.Err, "While writing interface:")
	return s.Rets()
}

func (i *Interface) RegisterImports(im *Imports) {
	//TODO: this is temporary, make it better
	for _, m := range i.methods {
		m.funcSig.RegisterImports(im)
	}
}

// PackageRef for the package Interface is in, fulfills Type interface.
func (i *Interface) PackageRef() PackageRef { return nil }

// Kind returns InterfaceKind, fulfills Type interface.
func (i *Interface) Kind() Kind { return InterfaceKind }

type interfaceMethod struct {
	funcSig FuncSig
	ifc     *Interface
}

func (im *interfaceMethod) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	s := gothicio.NewSumWriter(w)
	s.WriteString(im.funcSig.Name)
	s.WriteString("(")
	s.WriteString(typeSliceToString(im.funcSig.Args, pre, im.funcSig.Variadic))
	var end string
	if l := len(im.funcSig.Rets); l > 1 {
		s.WriteString(") (")
		end = ")"
	} else {
		s.WriteString(") ")
		end = ""
	}
	s.WriteString(typeSliceToString(im.funcSig.Rets, pre, false))
	s.WriteString(end)
	s.Err = errCtx(s.Err, "While writing interface method %s:", im.funcSig.Name)
	return s.Rets()
}

type interfaceRef struct {
	pkg  PackageRef
	name string
}

// DefInterface returns a reference to an interface in a package that fulfills
// Type.
func DefInterface(pkg PackageRef, name string) Type {
	return &interfaceRef{
		pkg:  pkg,
		name: name,
	}
}

func (i *interfaceRef) String() string { return i.pkg.Name() + "." + i.name }
func (i *interfaceRef) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString(pre.Prefix(i.pkg))
	sw.WriteString(i.name)
	sw.Err = errCtx(sw.Err, "While writing external interface reference %s", i.name)
	return sw.Rets()
}
func (i *interfaceRef) PackageRef() PackageRef { return i.pkg }
func (i *interfaceRef) Kind() Kind             { return InterfaceKind }

func (i *interfaceRef) RegisterImports(im *Imports) {
	im.AddRefImports(i.pkg)
}
