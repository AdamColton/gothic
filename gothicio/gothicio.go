package gothicio

import (
	"bytes"
	"io"
)

// StringWriter is the interface that wraps the WriteString method
type StringWriter interface {
	WriteString(string) (int, error)
}

// SumWriter is helper that wraps a Writer and sums the bytes written. If it
// encounters an error, it will stop writing.
type SumWriter struct {
	io.Writer
	Sum int64
	Err error
}

// NewSumWriter takes a Writer and returns a SumWriter
func NewSumWriter(w io.Writer) *SumWriter {
	return &SumWriter{Writer: w}
}

// WriteString writes a string to underlying Writer
func (s *SumWriter) WriteString(str string) (int, error) {
	return s.Write([]byte(str))
}

// WriteRune writes a rune to underlying Writer
func (s *SumWriter) WriteRune(r rune) { s.Write([]byte(string(r))) }

// Write fulfills io.Write
func (s *SumWriter) Write(b []byte) (int, error) {
	if s.Err != nil {
		return 0, s.Err
	}
	var n int
	n, s.Err = s.Writer.Write(b)
	s.Sum += int64(n)
	return n, s.Err
}

// MultiWrite takes a Writer and a slice of WriterTos and passes the Writer into
// each of them and writes the seperator between each.
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

// StringWriterTo fulfils the WriterTo interface and writes the string to the
// writer
type StringWriterTo string

// WriteTo fulfills the WriterTo interface.
func (s StringWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(s))
	return int64(n), err
}

// WriterToMerge merges multiple WriterTo into a single one. It tries to cast
// the first to an instance of WriterTos, making it efficient to merge multiple
// with successive calls.
func WriterToMerge(wts ...io.WriterTo) io.WriterTo {
	var w WriterTos
	if len(wts) == 0 {
		return w
	}
	if c, ok := wts[0].(WriterTos); ok {
		w = c
		wts = wts[1:]
	}
	for _, wt := range wts {
		if wt == nil {
			continue
		}
		w = append(w, wt)
	}
	return w
}

// WriterTos takes a slice of WriterTos and provides a single WriteTo method
// that will call each of them.
type WriterTos []io.WriterTo

// WriteTo invokes each of the WriterTos.
func (wts WriterTos) WriteTo(w io.Writer) (int64, error) {
	var sum int64
	for _, wt := range wts {
		n, err := wt.WriteTo(w)
		if err != nil {
			return int64(sum + n), err
		}
		sum += n
	}
	return int64(sum), nil
}

// TemplateExecutor is an interface representing the ExecuteTemplate method on
// a template.
type TemplateExecutor interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
	Execute(io.Writer, interface{}) error
}

type TemplateWrapper struct {
	TemplateExecutor
}

func (t TemplateWrapper) TemplateTo(name string, data interface{}) *TemplateTo {
	return &TemplateTo{
		TemplateExecutor: t,
		Name:             name,
		Data:             data,
	}
}

// TemplateTo writes a template and fulfils WriterTo. If Name is blank, the base
// template is used, otherwise the named template is used.
type TemplateTo struct {
	TemplateExecutor
	Name string
	Data interface{}
}

// NewTemplateTo returns a TemplateTo which fulfils WriterTo
func NewTemplateTo(template TemplateExecutor, name string, data interface{}) *TemplateTo {
	return &TemplateTo{
		TemplateExecutor: template,
		Name:             name,
		Data:             data,
	}
}

// WriteTo writes a template and fulfils WriterTo.
func (t *TemplateTo) WriteTo(w io.Writer) (int64, error) {
	var buf *bytes.Buffer
	var err error
	if Pool != nil {
		buf = Pool.Get()
	} else {
		buf = &bytes.Buffer{}
	}

	if t.Name == "" {
		err = t.Execute(buf, t.Data)
	} else {
		err = t.ExecuteTemplate(buf, t.Name, t.Data)
	}
	if err != nil {
		return 0, err
	}

	n, err := w.Write(buf.Bytes())
	if Pool != nil {
		Pool.Put(buf)
	}
	return int64(n), err
}

// BufferPool Gets and Puts Buffers, presumably backed by sync.Pool
type BufferPool interface {
	Get() *bytes.Buffer
	Put(buf *bytes.Buffer)
}

// Pool will be used if not nil when a Buffer is needed.
var Pool BufferPool

type BufferCloser struct {
	*bytes.Buffer
}

// Close allows BufferCloser to fill the closer interface
func (bc BufferCloser) Close() error {
	return nil
}
