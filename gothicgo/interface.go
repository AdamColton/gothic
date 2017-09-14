package gothicgo

import (
	"fmt"
	"strings"
)

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

// Adds a new Struct to an existing file
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

func (i *Interface) AddMethod(name string, args []Type, returns []Type, variadic bool) {
	i.methods = append(i.methods, &interfaceMethod{
		name:     name,
		args:     args,
		rets:     returns,
		variadic: variadic,
	})
}

func (i *Interface) Prepare() error { return nil }
func (i *Interface) Generate() error {
	i.file.AddCode(i.str())
	return nil
}

func (i *Interface) str() string {
	out := "type " + i.name + " interface{\n"
	for _, im := range i.methods {
		out += "\t" + im.str(i.file.Imports) + "\n"
	}
	out += "}\n\n"
	return out
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

func (i *Interface) Name() string { return i.name }
func (i *Interface) String() string {
	pkg := ""
	if i.file != nil && i.file.pkg != nil {
		pkg = i.file.pkg.Name + "."
	}
	return pkg + i.name
}
func (i *Interface) RelStr(pkg string) string {
	ipkg := ""
	if i.file != nil && i.file.pkg != nil && i.file.pkg.Name != pkg {
		ipkg = i.file.pkg.Name + "."
	}
	return ipkg + i.name
}
func (i *Interface) PackageName() string {
	if i.file != nil && i.file.pkg != nil {
		return i.file.pkg.Name
	}
	return ""
}
func (i *Interface) File() *File {
	return i.File()
}
func (i *Interface) Kind() Kind {
	return InterfaceKind
}
