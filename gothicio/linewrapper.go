package gothicio

import (
	"io"
	"unicode"
	"unicode/utf8"
)

var WrapWidth = 80

// LineWrapperContext returns contextual information to guide the writing
// operation. It is assumed that these values will not change with respect to an
// instance.
//
// LineWrapperContext only handle Unix style line endings.
type LineWrapperContext interface {
	WrapWidth() int
	Padding() string
}

// DefaulLineWrapperContext will return the package level WrapWidth and an empty
// string for the padding
type DefaulLineWrapperContext struct{}

// WrapWidth returns package level WrapWidth
func (DefaulLineWrapperContext) WrapWidth() int { return WrapWidth }

// Padding returns an empty string
func (DefaulLineWrapperContext) Padding() string { return "" }

// LineWrappingWriter
type LineWrappingWriter struct {
	LineWrapperContext
	*SumWriter
	sw         StringWriter
	padding    []byte
	onNewLine  bool
	start      int
	lineLength int
	lnPad      int
}

type LineWrapperContextWriter struct {
	io.Writer
	Width int
	Pad   string
}

func (lwcw LineWrapperContextWriter) WrapWidth() int  { return lwcw.Width }
func (lwcw LineWrapperContextWriter) Padding() string { return lwcw.Pad }

func NewLineWrappingWriter(w io.Writer) *LineWrappingWriter {
	sw, ok := w.(*SumWriter)
	if !ok {
		sw = &SumWriter{
			Writer: w,
		}
	} else {
		w = sw.Writer
	}

	lwc, ok := w.(LineWrapperContext)
	if !ok {
		lwc = DefaulLineWrapperContext{}
	}

	return &LineWrappingWriter{
		LineWrapperContext: lwc,
		SumWriter:          sw,
	}
}

func (w *LineWrappingWriter) setPadding() {
	w.padding = []byte(w.Padding())
	w.lnPad = utf8.RuneCount(w.padding)
}

func (w *LineWrappingWriter) Write(b []byte) (int, error) {
	if w.Err != nil {
		return 0, w.Err
	}
	if w.padding == nil {
		w.setPadding()
	}

	ww := w.WrapWidth()
	s0 := w.Sum

	start := 0
	lineLen := w.lineLength
	lastWS := -1
	done := true
	i := 0
	for i < len(b) {
		r, size := utf8.DecodeRune(b[i:])
		if r == '\n' {
			w.SumWriter.Write(b[start:i])
			w.SumWriter.Write(w.padding)
			w.onNewLine = true
			lineLen = w.lnPad
			i += size
			start = i // skip \n
			done = true
			continue
		}

		// 0xA0 is non-breaking space
		if unicode.IsSpace(r) && r != 0xA0 {
			lastWS = i
		} else {
			done = false
		}
		lineLen++
		if lineLen > ww && lastWS > 0 {
			w.SumWriter.Write(b[start:lastWS])
			start = lastWS + 1
			lastWS = -1
			w.WriteNewline()
			lineLen = w.lnPad
		}
		i += size
	}
	if !done {
		rest := b[start:]
		w.lineLength = len(rest)
		w.SumWriter.Write([]byte(rest))
	}
	return int(w.Sum - s0), w.Err
}

var nl = []byte("\n")

func (w *LineWrappingWriter) WriteNewline() (int, error) {
	if w.padding == nil {
		w.setPadding()
	}
	s0 := w.Sum
	w.SumWriter.Write(nl)
	w.SumWriter.Write(w.padding)
	w.lineLength = 0
	return int(w.Sum - s0), w.Err
}

func (w *LineWrappingWriter) WritePadding() (int, error) {
	if w.padding == nil {
		w.setPadding()
	}
	n, err := w.SumWriter.Write(w.padding)
	w.lineLength += w.lnPad
	return n, err
}
