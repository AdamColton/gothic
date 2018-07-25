package gothicgo

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

type ConstIotaBlock struct {
	t       Type
	rows    []string
	Comment string
	Iota    string
	file    interface {
		AddWriterTo(io.WriterTo) error
		Package() *Package
		Prefixer
	}
}

const (
	ErrEmptyConstIotaBlock = errStr("ConstIotaBlock requires at least one row")
	ErrConstIotaBlockType  = errStr("ConstIotaBlock requires a type")
)

type constRow string

func (c constRow) ScopeName() string {
	return string(c)
}

func (cib *ConstIotaBlock) WriteTo(w io.Writer) (int64, error) {
	if len(cib.rows) == 0 {
		return 0, ErrEmptyConstIotaBlock
	}
	if cib.t == nil {
		return 0, ErrConstIotaBlockType
	}
	s := gothicio.NewSumWriter(w)
	if cib.Comment != "" {
		NewComment(cib.Comment).WriteTo(s)
	}
	s.WriteString("const (\n\t")
	s.WriteString(cib.rows[0])
	s.WriteRune(' ')
	cib.t.PrefixWriteTo(s, cib.file)
	s.WriteString(" = ")
	if cib.Iota == "" {
		s.WriteString("iota")
	} else {
		s.WriteString(cib.Iota)
	}
	for _, r := range cib.rows[1:] {
		s.WriteString("\n\t")
		s.WriteString(r)
	}
	s.WriteString("\n)\n")
	if s.Err != nil {
		s.Err = errCtx(s.Err, "While writing ConstIotaBlock")
	}
	return s.Sum, s.Err
}

func (cib *ConstIotaBlock) Append(rows ...string) {
	pkg := cib.file.Package()
	for _, r := range rows {
		pkg.AddNamer(constRow(r))
	}
	cib.rows = append(cib.rows, rows...)
}

func (f *File) ConstIotaBlock(t Type, rows ...string) *ConstIotaBlock {
	cib := &ConstIotaBlock{
		t:    t,
		file: f,
		rows: rows,
	}
	f.AddWriterTo(cib)
	pkg := f.Package()
	for _, r := range rows {
		pkg.AddNamer(constRow(r))
	}
	return cib
}
