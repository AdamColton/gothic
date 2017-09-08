package gothichtml

import (
	"io"
)

type writer struct {
	sw            StringWriter
	onNewLine     bool
	start         int
	padding       string
	parentPadding string
}

func newWriter(w io.Writer) *writer {
	return &writer{
		sw: ToStringWriter(w),
	}
}

func (w *writer) write(str string) {
	w.onNewLine = false
	w.sw.WriteString(str)
}

func (w *writer) nl() {
	w.onNewLine = true
	w.sw.WriteString(NewLine)
	w.sw.WriteString(w.padding)
}

func (w *writer) pnl() {
	w.onNewLine = true
	w.sw.WriteString(NewLine)
	w.sw.WriteString(w.parentPadding)
}

func (w *writer) inc() *writer {
	cp := *w
	cp.parentPadding = cp.padding
	cp.padding += Padding
	return &cp
}

// Node is any node that can be in an html document
type Node interface {
	Parent() Node
	Write(io.Writer)
	write(*writer)
	setParent(Node)
}

// ContainerNode is a Node that has child nodes
type ContainerNode interface {
	Node
	Children() int
	Child(int) Node
	AddChildren(...Node)
	RemoveChild(int)
	frag() *fragment
}

// TagNode is a node with a Tag Name and has attributes.
type TagNode interface {
	Node
	Name() string
	Attributes() []string
	Attribute(string) (string, bool)
	AddAttributes(attrs ...string)
	Remove(key string)
}
