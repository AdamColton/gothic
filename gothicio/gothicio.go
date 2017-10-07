package gothicio

import (
	"io"
)

type SumWriter struct {
	io.Writer
	Sum int64
	Err error
}

func NewSumWriter(w io.Writer) *SumWriter {
	return &SumWriter{Writer: w}
}

func (s *SumWriter) WriteString(str string) { s.Write([]byte(str)) }
func (s *SumWriter) WriteRune(r rune)       { s.Write([]byte(string(r))) }

func (s *SumWriter) Write(b []byte) {
	if s.Err != nil {
		return
	}
	var n int
	n, s.Err = s.Writer.Write(b)
	s.Sum += int64(n)
}

func (s *SumWriter) Wrap(wt io.WriterTo) {
	if s.Err != nil {
		return
	}
	var n int64
	n, s.Err = wt.WriteTo(s.Writer)
	s.Sum += n
}

func MultiWrite(w io.Writer, tos ...io.WriterTo) (int64, error) {
	var s int64
	for _, t := range tos {
		n, err := t.WriteTo(w)
		if err != nil {
			return s, err
		}
		s += int64(n)
	}
	return s, nil
}
