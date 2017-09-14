package gothicgo

import (
	"strings"
)

// Func represents a Go function
type Func struct {
	*Imports
	name     string
	Args     []*NameType
	Rets     []*NameType
	Body     string
	Variadic bool
	//TODO: now that Imports is shared with file, this shouldn't be exported.
	File *File
}

// NewFunc takes the function name and arguments and returns a Func
func NewFunc(name string, args ...*NameType) *Func {
	return &Func{
		Imports: NewImports(pkgBuiltin),
		name:    name,
		Args:    args,
	}
}

func (f *File) NewFunc(name string, args ...*NameType) *Func {
	fn := &Func{
		Imports: f.Imports,
		name:    name,
		Args:    args,
		File:    f,
	}
	f.AddGenerators(fn)
	return fn
}

func (f *Func) GetName() string     { return f.name }
func (f *Func) SetName(name string) { f.name = name }

// Returns sets the returns on a Func
func (f *Func) Returns(rets ...*NameType) { f.Rets = rets }

func nameTypeSliceToString(imp *Imports, nts []*NameType, variadic bool) string {
	l := len(nts)
	var s = make([]string, l)
	l--
	for i := l; i >= 0; i-- {
		if ts := nts[i].T.RelStr(imp); i < l && ts == nts[i+1].T.RelStr(imp) {
			s[i] = nts[i].N
		} else if i == l && variadic {
			s[i] = nts[i].N + " ..." + ts
		} else {
			s[i] = nts[i].N + " " + ts
		}
	}
	return strings.Join(s, ", ")
}

func (f *Func) RelSignature(pkg string) string {
	s := make([]string, 6)
	s[0] = f.name
	s[1] = "("
	s[2] = nameTypeSliceToString(f.Imports, f.Args, f.Variadic)
	if l := len(f.Rets); l > 1 || (l == 1 && f.Rets[0].N != "") {
		s[3] = ") ("
		s[5] = ")"
	} else {
		s[3] = ")"
		s[5] = ""
	}
	s[4] = nameTypeSliceToString(f.Imports, f.Rets, false)

	return strings.Join(s, "")
}

// String outputs the entire function as a string
func (f *Func) String() string {
	s := make([]string, 9)
	s[0] = "func "
	s[1] = f.name
	s[2] = "("
	s[3] = nameTypeSliceToString(f.Imports, f.Args, f.Variadic)
	if l := len(f.Rets); l > 1 || (l == 1 && f.Rets[0].N != "") {
		s[4] = ") ("
		s[6] = ") {\n"
	} else {
		s[4] = ")"
		s[6] = " {\n"
	}
	s[5] = nameTypeSliceToString(f.Imports, f.Rets, false)
	s[7] = f.Body
	s[8] = "\n}\n\n"
	return strings.Join(s, "")
}

func (f *Func) Prepare() error {
	return nil
}

func (f *Func) Generate() error {
	f.File.AddCode(f.String())
	return nil
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
	return f.RelStr(nil)
}

func (f *fnT) RelStr(imp *Imports) string {
	s := make([]string, 5)
	s[0] = "func("
	args := make([]string, len(f.fn.Args))
	for i, arg := range f.fn.Args {
		args[i] = arg.T.RelStr(imp)
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
		rets[i] = ret.T.RelStr(imp)
	}
	s[3] = strings.Join(rets, ", ")
	return strings.Join(s, "")
}

func (f *fnT) String() string {
	return f.RelStr(nil)
}

func (f *fnT) PackageRef() PackageRef {
	if f.fn.File != nil {
		return f.fn.File.pkg.Ref
	}
	return pkgBuiltin
}

func (f *fnT) File() *File       { return f.File() }
func (f *fnT) Kind() Kind        { return FuncKind }
func (f *fnT) Args() []*NameType { return f.fn.Args }
func (f *fnT) Rets() []*NameType { return f.fn.Rets }
func (f *fnT) Variadic() bool    { return f.fn.Variadic }
