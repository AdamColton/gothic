package gothicio

import (
	"io"
)

type StringWriter interface {
	WriteString(string) (int, error)
}

type SumWriter struct {
	io.Writer
	Sum int64
	Err error
}

func NewSumWriter(w io.Writer) *SumWriter {
	return &SumWriter{Writer: w}
}

func (s *SumWriter) WriteString(str string) (int, error) {
	return s.Write([]byte(str))
}
func (s *SumWriter) WriteRune(r rune) { s.Write([]byte(string(r))) }

func (s *SumWriter) Write(b []byte) (int, error) {
	if s.Err != nil {
		return 0, s.Err
	}
	var n int
	n, s.Err = s.Writer.Write(b)
	s.Sum += int64(n)
	return n, s.Err
}

func MultiWrite(w io.Writer, tos []io.WriterTo, seperator string) (int64, error) {
	sbs := []byte(seperator)
	var s int64
	for i, t := range tos {
		if i != 0 {
			n, err := w.Write(sbs)
			if err != nil {
				return s, err
			}
			s += int64(n)
		}
		n, err := t.WriteTo(w)
		if err != nil {
			return s, err
		}
		s += int64(n)
	}
	return s, nil
}
