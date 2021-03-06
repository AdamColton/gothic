package gothicgo

import (
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"strings"
)

// InterfaceEmbeddable allows one interface to be embedded in another
type InterfaceEmbeddable interface {
	Type
	InterfaceEmbedName() string
}

// Interface is used to generate an interface
type Interface struct {
	methods  []*interfaceMethod
	embedded []InterfaceEmbeddable
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

func (i *Interface) Embed(embed InterfaceEmbeddable) {
	i.embedded = append(i.embedded, embed)
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
	if len(i.methods) == 0 && len(i.embedded) == 0 {
		n, err := w.Write([]byte("interface{}"))
		return int64(n), err
	}
	s := gothicio.NewSumWriter(w)
	s.WriteString("interface {")
	for _, e := range i.embedded {
		s.WriteString("\n\t")
		e.PrefixWriteTo(s, pre)
	}
	for _, m := range i.methods {
		s.WriteString("\n\t")
		m.PrefixWriteTo(s, pre)
	}
	s.WriteString("\n}")
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
