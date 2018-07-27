package gothicgo

import (
	"fmt"
	"github.com/adamcolton/gothic/bufpool"
	"github.com/adamcolton/gothic/gothicio"
	"io"
	"strings"
)

var _ Type = &FuncSig{}

type FuncSig struct {
	Name     string
	Args     []NameType
	Rets     []NameType
	Variadic bool
}

func NewFuncSig(name string, args ...NameType) FuncSig {
	return FuncSig{
		Name: name,
		Args: args,
	}
}

func (f FuncSig) Type() Type {
	args := make([]NameType, len(f.Args))
	rets := make([]NameType, len(f.Rets))
	for i, a := range f.Args {
		args[i].T = a.T
	}
	for i, r := range f.Rets {
		rets[i].T = r.T
	}
	return FuncSig{
		Name:     f.Name,
		Args:     args,
		Rets:     rets,
		Variadic: f.Variadic,
	}
}

func (f FuncSig) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := gothicio.NewSumWriter(w)
	sw.WriteString("func")
	if f.Name != "" {
		sw.WriteRune(' ')
		sw.WriteString(f.Name)
	}
	sw.WriteRune('(')
	var str string
	str, sw.Err = nameTypeSliceToString(pre, f.Args, f.Variadic)
	sw.WriteString(str)
	end := ""
	if len(f.Rets) > 1 {
		sw.WriteString(") (")
		end = ")"
	} else {
		sw.WriteString(") ")
	}
	str, sw.Err = nameTypeSliceToString(pre, f.Rets, false)
	sw.WriteString(str)
	sw.WriteString(end)

	return sw.Sum, sw.Err
}

func (f FuncSig) Kind() Kind             { return FuncKind }
func (f FuncSig) PackageRef() PackageRef { return nil }
func (f FuncSig) RegisterImports(i *Imports) {
	for _, arg := range f.Args {
		arg.T.RegisterImports(i)
	}
	for _, ret := range f.Rets {
		ret.T.RegisterImports(i)
	}
}
func (f FuncSig) String() string { return typeToString(f, DefaultPrefixer) }

func nameTypeSliceToString(pre Prefixer, nts []NameType, variadic bool) (string, error) {
	if len(nts) == 0 {
		return "", nil
	}
	if nts[0].N == "" {
		return unnamedTypeSliceToString(pre, nts, variadic)
	}
	return namedTypeSliceToString(pre, nts, variadic)
}

func unnamedTypeSliceToString(pre Prefixer, nts []NameType, variadic bool) (string, error) {
	l := len(nts)
	var s = make([]string, l)
	l--
	b := bufpool.Get()
	for i := l; i >= 0; i-- {
		if nts[i].N != "" {
			return "", errStr("mixed named and unnamed function parameters")
		}
		nts[i].T.PrefixWriteTo(b, pre)
		if i == l && variadic {
			s[i] = fmt.Sprintf("...%s", b.String())
		} else if nts[i].N != "" {
			s[i] = fmt.Sprint(b.String())
		} else {
			s[i] = b.String()
		}
		b.Reset()
	}
	bufpool.Put(b)
	return strings.Join(s, ", "), nil
}

func namedTypeSliceToString(pre Prefixer, nts []NameType, variadic bool) (string, error) {
	l := len(nts)
	var s = make([]string, l)
	l--
	b1 := bufpool.Get()
	b2 := bufpool.Get()
	for i := l; i >= 0; i-- {
		if nts[i].N == "" {
			return "", errStr("mixed named and unnamed function parameters")
		}
		nts[i].T.PrefixWriteTo(b1, pre)
		if i < l {
			nts[i+1].T.PrefixWriteTo(b2, pre)
		}
		if i < l && b1.String() == b2.String() {
			s[i] = nts[i].N
		} else if i == l && variadic {
			s[i] = fmt.Sprintf("%s ...%s", nts[i].N, b1.String())
		} else if nts[i].N != "" {
			s[i] = fmt.Sprintf("%s %s", nts[i].N, b1.String())
		} else {
			s[i] = b1.String()
		}
		b1.Reset()
		b2.Reset()
	}
	bufpool.Put(b1)
	bufpool.Put(b2)
	return strings.Join(s, ", "), nil
}

// Returns sets the return types on the function
func (f *FuncSig) Returns(rets ...NameType) {
	f.Rets = rets
}

// Func function written to a Go file
type Func struct {
	FuncSig
	Body PrefixWriterTo
	//TODO: now that Imports is shared with file, this shouldn't be exported.
	Comment string
	file    *File
}

// NewFunc returns a new Func with File set and add the function to file's
// generators so that when the file is generated, the func will be generated as
// part of the file.
func (f *File) NewFunc(name string, args ...NameType) (*Func, error) {
	fn := &Func{
		FuncSig: NewFuncSig(name, args...),
		file:    f,
	}
	return fn, errCtx(f.AddWriterTo(fn), "Adding func %s to file %s", name, f.name)
}

func (f *Func) ScopeName() string { return f.Name }

func (f *Func) WriteTo(w io.Writer) (int64, error) {
	return f.PrefixWriteTo(w, f.file)
}

// Prepare adds all the types used in the Args and Rets to the file import.
func (f *Func) Prepare() error {
	f.RegisterImports(f.file.Imports)
	if ri, ok := f.Body.(RegisterImports); ok {
		ri.RegisterImports(f.file.Imports)
	}
	return nil
}

// UnnamedReturns set the return types on the function, all with no names.
func (f *Func) UnnamedReturns(rets ...Type) {
	f.Rets = Rets(rets...)
}

// String outputs the entire function as a string. It will invoke f.Body, so it
// is best used for testing or else it can be called out of order in the
// Prepare/Generate cycle.
func (f *Func) String() string {
	buf := bufpool.Get()
	f.WriteTo(buf)
	return bufpool.PutStr(buf)
}

// WriteTo writes the Func to a writer
func (f *Func) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	s := gothicio.NewSumWriter(w)
	if f.Comment != "" {
		NewComment(strings.Join([]string{f.Name, f.Comment}, " ")).WriteTo(s)
	}
	f.FuncSig.PrefixWriteTo(w, pre)
	s.WriteString(" {\n")
	if f.Body != nil {
		f.Body.PrefixWriteTo(s, pre)
	}
	s.WriteString("\n}")
	s.Err = errCtx(s.Err, "While writing func %s", f.Name)
	return s.Sum, s.Err
}

func (f *Func) BodyWriterTo(w io.WriterTo) {
	f.Body = IgnorePrefixer{w}
}

func (f *Func) BodyString(str string) {
	f.Body = IgnorePrefixer{gothicio.StringWriterTo(str)}
}

// Call produces a invocation of the function and fulfills the FuncCaller
// interface
func (f *Func) Call(pre Prefixer, args ...string) string {
	buf := bufpool.Get()
	buf.WriteString(pre.Prefix(f.file.Package()))
	buf.WriteString(f.Name)
	buf.WriteRune('(')
	buf.WriteString(strings.Join(args, ", "))
	buf.WriteRune(')')
	str := buf.String()
	bufpool.Put(buf)
	return str
}
