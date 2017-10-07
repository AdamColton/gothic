package gothicgo

import (
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"strings"
)

// Func represents a Go function.
type Func struct {
	*Imports
	name     string
	Args     []*NameType
	Rets     []*NameType
	Body     func() (string, error)
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

// NewFunc returns a new Func with File set and add the function to file's
// generators so that when the file is generated, the func will be generated as
// part of the file.
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

// Name of the function
func (f *Func) Name() string { return f.name }

// SetName of the function
func (f *Func) SetName(name string) { f.name = name }

// Returns sets the returns on a Func
func (f *Func) Returns(rets ...*NameType) { f.Rets = rets }

func nameTypeSliceToString(pre Prefixer, nts []*NameType, variadic bool) string {
	l := len(nts)
	var s = make([]string, l)
	l--
	for i := l; i >= 0; i-- {
		if ts := nts[i].T.RelStr(pre); i < l && ts == nts[i+1].T.RelStr(pre) {
			s[i] = nts[i].N
		} else if i == l && variadic {
			s[i] = nts[i].N + " ..." + ts
		} else if nts[i].N != "" {
			s[i] = nts[i].N + " " + ts
		} else {
			s[i] = ts
		}
	}
	return strings.Join(s, ", ")
}

// String outputs the entire function as a string. It will invoke f.Body, so it
// is best used for testing or else it can be called out of order in the
// Prepare/Generate cycle.
func (f *Func) String() string {
	buf := bufpool.Get()
	f.WriteTo(buf)
	return bufpool.PutStr(buf)
}

func (f *Func) WriteTo(w io.Writer) (int64, error) {
	s := gothicio.NewSumWriter(w)
	body, err := f.Body()
	if err != nil {
		return 0, err
	}
	s.WriteString("func ")
	s.WriteString(f.name)
	writeArgsRets(s, f.Imports, f.Args, f.Rets, f.Variadic)
	s.WriteString(" {\n")
	s.WriteString(body)
	s.WriteString("\n}\n\n")
	return s.Sum, s.Err
}

// Prepare adds all the types used in the Args and Rets to the file import.
func (f *Func) Prepare() error {
	if f.File != nil {
		for _, arg := range f.Args {
			f.File.AddRefImports(arg.Type().PackageRef())
		}
		for _, ret := range f.Rets {
			f.File.AddRefImports(ret.Type().PackageRef())
		}
	}
	return nil
}

// Generate writes the function to the file
func (f *Func) Generate() error {
	f.File.AddWriteTo(f)
	return nil
}

// Type returns a FuncType which fulfills the Type interface
func (f *Func) Type() FuncType { return &fnT{f} }

// Call produces a invocation of the function and fulfills the FuncCaller
// interface
func (f *Func) Call(pre Prefixer, args ...string) string {
	buf := bufpool.Get()
	buf.WriteString(pre.Prefix(f.File.Package()))
	buf.WriteString(f.name)
	buf.WriteRune('(')
	buf.WriteString(strings.Join(args, ", "))
	buf.WriteRune(')')
	str := buf.String()
	bufpool.Put(buf)
	return str
}

// FuncType fulfills Type and adds additional information
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
	return f.RelStr(DefaultPrefixer)
}

func (f *fnT) RelStr(pre Prefixer) string {
	buf := bufpool.Get()
	buf.WriteString("func(")
	for i, arg := range f.fn.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.T.RelStr(pre))
	}
	lr := len(f.fn.Rets)
	if lr > 1 {
		buf.WriteString(") (")
	} else {
		buf.WriteString(") ")
	}
	for i, ret := range f.fn.Rets {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(ret.T.RelStr(pre))
	}

	if lr > 1 {
		buf.WriteString(")")
	}
	return bufpool.PutStr(buf)
}

func (f *fnT) String() string {
	return f.RelStr(DefaultPrefixer)
}

func (f *fnT) PackageRef() PackageRef {
	if f.fn.File != nil {
		return f.fn.File.pkg
	}
	return pkgBuiltin
}

func (f *fnT) File() *File       { return f.File() }
func (f *fnT) Kind() Kind        { return FuncKind }
func (f *fnT) Args() []*NameType { return f.fn.Args }
func (f *fnT) Rets() []*NameType { return f.fn.Rets }
func (f *fnT) Variadic() bool    { return f.fn.Variadic }

type funcCall struct {
	pkg  PackageRef
	name string
}

// FuncCaller produces a string that will call a function. It handles the
// correct prefixing of the function call.
type FuncCaller interface {
	Call(Prefixer, ...string) string
}

// FuncCall defines a callable function in another package
func FuncCall(pkg PackageRef, name string) FuncCaller {
	return &funcCall{
		pkg:  pkg,
		name: name,
	}
}

func (f *funcCall) Call(pre Prefixer, args ...string) string {
	buf := bufpool.Get()
	buf.WriteString(pre.Prefix(f.pkg))
	buf.WriteString(f.name)
	buf.WriteRune('(')
	buf.WriteString(strings.Join(args, ", "))
	buf.WriteRune(')')
	str := buf.String()
	bufpool.Put(buf)
	return str
}

func writeArgsRets(s *gothicio.SumWriter, pre Prefixer, args, rets []*NameType, variadic bool) {
	s.WriteRune('(')
	s.WriteString(nameTypeSliceToString(pre, args, variadic))
	if l := len(rets); l > 1 || (l == 1 && rets[0].N != "") {
		s.WriteString(") (")
		s.WriteString(nameTypeSliceToString(pre, rets, false))
		s.WriteRune(')')
	} else {
		s.WriteString(") ")
		if l == 1 {
			s.WriteString(nameTypeSliceToString(pre, rets, false))
		}
	}
}
