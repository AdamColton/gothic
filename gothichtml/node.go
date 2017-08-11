package gothichtml

import (
	"io"
)

// Node is any node that can be in an html document
type Node interface {
	Parent() Node
	Write(io.Writer)
	write(StringWriter, string)
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
