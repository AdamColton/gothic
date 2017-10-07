package gothicgo

import (
	"fmt"
	"io"
	"strings"
)

// Interface is used to generate an interface
type Interface struct {
	name    string
	file    *File
	methods []*interfaceMethod
}

type interfaceMethod struct {
	name     string
	args     []Type
	rets     []Type
	variadic bool
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
	})
}

// Prepare fulfills the generator interface
func (i *Interface) Prepare() error { return nil }

// Generate adds the interface to the file and fulfills the generator interface
func (i *Interface) Generate() error {
	i.file.AddWriteTo(i)
	return nil
}

func (i *Interface) WriteTo(w io.Writer) (int64, error) {
	s := SumWriter{W: w}
	s.WriteString("type ")
	s.WriteString(i.Name())
	s.WriteString(" interface{\n")
	for _, im := range i.methods {
		s.WriteString("\t")
		s.WriteString(im.str(i.file.Imports))
		s.WriteString("\n")
	}
	s.WriteString("}\n\n")
	return s.Sum, s.Err
}

func typeSliceToString(ts []Type, imp *Imports, variadic bool) string {
	l := len(ts)
	var s = make([]string, l)
	l--
	for i, t := range ts {
		if i == l && variadic {
			s[i] = " ..." + t.RelStr(imp)
		} else {
			s[i] = t.RelStr(imp)
		}
	}
	return strings.Join(s, ", ")
}

func (im *interfaceMethod) str(imp *Imports) string {
	s := make([]string, 6)
	s[0] = im.name
	s[1] = "("
	s[2] = typeSliceToString(im.args, imp, im.variadic)
	if l := len(im.rets); l > 1 {
		s[3] = ") ("
		s[5] = ")"
	} else {
		s[3] = ") "
		s[5] = ""
	}
	s[4] = typeSliceToString(im.rets, imp, false)

	return strings.Join(s, "")
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

// RelStr returns a string with the interface name and package if necessary.
func (i *Interface) RelStr(pre Prefixer) string {
	return pre.Prefix(i.file.pkg) + i.name
}

// PackageRef for the package Interface is in, fulfills Type interface.
func (i *Interface) PackageRef() PackageRef { return i.file.pkg }

// File that the interface is in, fulfills Type interface.
func (i *Interface) File() *File { return i.File() }

// Kind returns InterfaceKind, fulfills Type interface.
func (i *Interface) Kind() Kind { return InterfaceKind }

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
func (i *interfaceRef) RelStr(pre Prefixer) string {
	return pre.Prefix(i.pkg) + i.name
}
func (i *interfaceRef) PackageRef() PackageRef { return i.pkg }
func (i *interfaceRef) File() *File            { return nil }
func (i *interfaceRef) Kind() Kind             { return InterfaceKind }
