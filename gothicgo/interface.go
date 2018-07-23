package gothicgo

import (
	"fmt"
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"strings"
)

// Interface is used to generate an interface
type Interface struct {
	name    string
	file    *File
	methods []io.WriterTo
}

// NewInterface adds a new interface to an existing file
func (f *File) NewInterface(name string) (*Interface, error) {
	if i, found := f.pkg.interfaces[name]; found {
		return i, fmt.Errorf("Duplicate Interface")
	}
	i := &Interface{
		name: name,
		file: f,
	}
	f.pkg.interfaces[name] = i
	f.AddGenerators(i)
	return i, nil
}

// AddMethod to the interface
func (i *Interface) AddMethod(name string, args []Type, returns []Type, variadic bool) {
	i.methods = append(i.methods, &interfaceMethod{
		name:     name,
		args:     args,
		rets:     returns,
		variadic: variadic,
		ifc:      i,
	})
}

// Prepare fulfills the generator interface
func (i *Interface) Prepare() error { return nil }

// Generate adds the interface to the file and fulfills the generator interface
func (i *Interface) Generate() error {
	i.file.AddWriterTo(i)
	return nil
}

// WriteTo writes the Interface code to the writer
func (i *Interface) WriteTo(w io.Writer) (int64, error) {
	s := gothicio.NewSumWriter(w)
	s.WriteString("type ")
	s.WriteString(i.Name())
	s.WriteString(" interface{")
	gothicio.MultiWrite(s, i.methods, "\n\t")
	if len(i.methods) > 0 {
		s.WriteString("\n}\n\n")
	} else {
		s.WriteString("}\n\n")
	}
	s.Err = errCtx(s.Err, "While writing interface %s:", i.name)
	return s.Sum, s.Err
}

func typeSliceToString(ts []Type, pre Prefixer, variadic bool) string {
	l := len(ts)
	var s = make([]string, l)
	l--
	buf := bufpool.Get()
	for i, t := range ts {
		t.PrefixWriteTo(buf, pre)
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

// Name gets the interface name and fulfills the Type interface
func (i *Interface) Name() string { return i.name }

// SetName allows the name of the interface to be changed
func (i *Interface) SetName(name string) { i.name = name }

// String returns the interface package and name and fulfills the Type interface
func (i *Interface) String() string {
	pkg := ""
	if i.file != nil && i.file.pkg != nil {
		pkg = i.file.pkg.name + "."
	}
	return pkg + i.name
}

// WriteTo writes the interface name and package if necessary.
func (i *Interface) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString(pre.Prefix(i.file.pkg))
	sw.WriteString(i.name)
	return sw.Rets()
}

// PackageRef for the package Interface is in, fulfills Type interface.
func (i *Interface) PackageRef() PackageRef { return i.file.pkg }

// File that the interface is in, fulfills Type interface.
func (i *Interface) File() *File { return i.File() }

// Kind returns InterfaceKind, fulfills Type interface.
func (i *Interface) Kind() Kind { return InterfaceKind }

type interfaceMethod struct {
	name     string
	args     []Type
	rets     []Type
	variadic bool
	ifc      *Interface
}

func (im *interfaceMethod) WriteTo(w io.Writer) (int64, error) {
	s := gothicio.NewSumWriter(w)
	s.WriteString(im.name)
	s.WriteString("(")
	s.WriteString(typeSliceToString(im.args, im.ifc.file.Imports, im.variadic))
	var end string
	if l := len(im.rets); l > 1 {
		s.WriteString(") (")
		end = ")"
	} else {
		s.WriteString(") ")
		end = ""
	}
	s.WriteString(typeSliceToString(im.rets, im.ifc.file.Imports, false))
	s.WriteString(end)
	s.Err = errCtx(s.Err, "While writing interface method %s:", im.name)
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

func (i *interfaceRef) Name() string   { return i.name }
func (i *interfaceRef) String() string { return i.pkg.Name() + "." + i.name }
func (i *interfaceRef) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString(pre.Prefix(i.pkg))
	sw.WriteString(i.name)
	return sw.Rets()
}
func (i *interfaceRef) PackageRef() PackageRef { return i.pkg }
func (i *interfaceRef) File() *File            { return nil }
func (i *interfaceRef) Kind() Kind             { return InterfaceKind }
