package gothicgo

import (
	"strings"
)

// Func represents a Go function
type Func struct {
	*Imports
	Name     string
	Args     []*NameType
	Rets     []*NameType
	Body     string
	Variadic bool
	File     *File
}

// NewFunc takes the function name and arguments and returns a Func
func NewFunc(name string, args ...*NameType) *Func {
	return &Func{
		Imports: NewImports(),
		Name:    name,
		Args:    args,
	}
}

func (f *File) NewFunc(name string, args ...*NameType) *Func {
	fn := &Func{
		Imports: NewImports(),
		Name:    name,
		Args:    args,
		File:    f,
	}
	f.AddFragGen(fn)
	return fn
}

// Returns sets the returns on a Func
func (f *Func) Returns(rets ...*NameType) { f.Rets = rets }

func nameTypeSliceToString(nts []*NameType, pkgName string, variadic bool) string {
	l := len(nts)
	var s = make([]string, l)
	l--
	for i := l; i >= 0; i-- {
		if ts := nts[i].T.RelStr(pkgName); i < l && ts == nts[i+1].T.RelStr(pkgName) {
			s[i] = nts[i].N
		} else if i == l && variadic {
			s[i] = nts[i].N + " ..." + ts
		} else {
			s[i] = nts[i].N + " " + ts
		}
	}
	return strings.Join(s, ", ")
}

// String outputs the entire function as a string
func (f *Func) String() string {
	pkgName := ""
	if f.File != nil {
		pkgName = f.File.Package().Name
	}

	s := make([]string, 9)
	s[0] = "func "
	s[1] = f.Name
	s[2] = "("
	s[3] = nameTypeSliceToString(f.Args, pkgName, f.Variadic)
	if l := len(f.Rets); l > 1 || (l == 1 && f.Rets[0].N != "") {
		s[4] = ") ("
		s[6] = ") {\n"
	} else {
		s[4] = ")"
		s[6] = " {\n"
	}
	s[5] = nameTypeSliceToString(f.Rets, pkgName, false)
	s[7] = f.Body
	s[8] = "\n}\n\n"
	return strings.Join(s, "")
}

func (f *Func) Prepare() {
	if f.File != nil {
		pkgName := f.File.Package().Name
		for _, a := range f.Args {
			if aPkgName := a.T.PackageName(); aPkgName != pkgName {
				f.File.AddPackageImport(aPkgName)
			}
		}
		for _, r := range f.Rets {
			if rPkgName := r.T.PackageName(); rPkgName != pkgName {
				f.File.AddPackageImport(rPkgName)
			}
		}
		f.File.Imports.AddImports(f.Imports)
	}
}

func (f *Func) Generate() []string {
	return []string{f.String()}
}

func (f *Func) Type() FuncType {
	return &fnT{f}
}

type FuncType interface {
	Type
	Args() []*NameType
	Rets() []*NameType
	Variadic() bool
}

type fnT struct {
	fn *Func
}

func (f *fnT) Name() string {
	return f.RelStr(f.PackageName())
}

func (f *fnT) RelStr(pkg string) string {
	s := make([]string, 5)
	s[0] = "func("
	args := make([]string, len(f.fn.Args))
	for i, arg := range f.fn.Args {
		args[i] = arg.T.RelStr(pkg)
	}
	s[1] = strings.Join(args, ", ")
	lr := len(f.fn.Rets)
	if lr > 1 {
		s[2] = ") ("
		s[4] = ")"
	} else {
		s[2] = ") "
	}
	rets := make([]string, lr)
	for i, ret := range f.fn.Rets {
		rets[i] = ret.T.RelStr(pkg)
	}
	s[3] = strings.Join(rets, ", ")
	return strings.Join(s, "")
}

func (f *fnT) String() string {
	return f.RelStr("")
}

func (f *fnT) PackageName() string {
	if f.fn.File != nil {
		return f.fn.File.pkg.Name
	}
	return ""
}

func (f *fnT) File() *File       { return f.File() }
func (f *fnT) Kind() Kind        { return FuncKind }
func (f *fnT) Args() []*NameType { return f.fn.Args }
func (f *fnT) Rets() []*NameType { return f.fn.Rets }
func (f *fnT) Variadic() bool    { return f.fn.Variadic }
