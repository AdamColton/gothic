package gothichtml

import (
	"bytes"
	"io"
	"strings"
)

// NewLine is the string that will be used when rendering HTML to a writer for a
// newline.
var NewLine = "\n"

// Padding is the string that is used to indent html
var Padding = "  "

// StringWriter takes a string as an argument (instead of []bytes). It is
// fulfilled by bytes.Buffer.
type StringWriter interface {
	WriteString(s string) (n int, err error)
}

// StringWriterWrapper can wrap any io.Writer to fulfill StringWriter
type StringWriterWrapper struct {
	io.Writer
}

// WriteString fulfils StringWriter on StringWriterWrapper. It just casts the
// string to []byte
func (w StringWriterWrapper) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

// ToStringWriter checks if the writer can be cast to a StringWriter and if it
// cannot, it converts it using a StringWriterWrapper.
func ToStringWriter(w io.Writer) StringWriter {
	if sw, ok := w.(StringWriter); ok {
		return sw
	}
	return StringWriterWrapper{w}
}

// String uses bytes.Buffer to render html as a string. This tends to be useful
// in testing, but less so in production code.
func String(node Node) string {
	var buf bytes.Buffer
	node.write(newWriter(&buf))
	return buf.String()
}

// Classes is a helper function that returns the classes on a TagNode as a slice
// of strings.
func Classes(node TagNode) []string {
	classes, _ := node.Attribute("class")
	return strings.Fields(classes)
}

// parent is a helper that is embeded to hold a reference to self and parent and
// fulfill the Parent() and setParent() methods on Node.
type parent struct {
	parent Node
	self   Node
}

func newParent(self Node) *parent {
	return &parent{self: self}
}

func (p *parent) Parent() Node {
	return p.parent
}

func (p *parent) setParent(newparent Node) {
	p.parent = newparent
}
