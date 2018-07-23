package gothicgo

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

type ConstIotaBlock struct {
	Type    Type
	Rows    []string
	Comment string
	Iota    string
	Prefixer
}

const (
	ErrEmptyConstIotaBlock = errStr("ConstIotaBlock requires at least one row")
	ErrConstIotaBlockType  = errStr("ConstIotaBlock requires a type")
)

func (cib *ConstIotaBlock) WriteTo(w io.Writer) (int64, error) {
	if len(cib.Rows) == 0 {
		return 0, ErrEmptyConstIotaBlock
	}
	if cib.Type == nil {
		return 0, ErrConstIotaBlockType
	}
	s := gothicio.NewSumWriter(w)
	if cib.Comment != "" {
		NewComment(cib.Comment).WriteTo(s)
	}
	s.WriteString("const (\n\t")
	s.WriteString(cib.Rows[0])
	s.WriteRune(' ')
	cib.Type.PrefixWriteTo(s, cib)
	s.WriteString(" = ")
	if cib.Iota == "" {
		s.WriteString("iota")
	} else {
		s.WriteString(cib.Iota)
	}
	for _, r := range cib.Rows[1:] {
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
	cib.Rows = append(cib.Rows, rows...)
}

func (f *File) ConstIotaBlock(t Type, rows ...string) *ConstIotaBlock {
	cib := &ConstIotaBlock{
		Type:     t,
		Prefixer: f,
		Rows:     rows,
	}
	f.AddWriterTo(cib)
	return cib
}
