package gothichtml

import (
	"io"
	"strings"
)

// Tag is any tag with an opening and closing that may have children
type Tag struct {
	tag string
	attributes
	*fragment
}

// NewTag creates a Tag. It requires a tag name and optionally takes attributes.
func NewTag(tag string, attrs ...string) *Tag {
	t := &Tag{
		tag:        tag,
		attributes: newAttributes(attrs),
		fragment:   &fragment{},
	}
	t.parent = newParent(t)
	return t
}

// Name returns the tag name
func (t *Tag) Name() string { return t.tag }

// Write a tag to a writer
func (t *Tag) Write(w io.Writer) {
	t.write(ToStringWriter(w), NewLine)
}

func (t *Tag) write(w StringWriter, padding string) {
	sw := w.WriteString
	sw("<")
	sw(t.tag)
	t.attributes.write(w)
	sw(">")

	multiline := true
	if l := len(t.fragment.children); l == 0 {
		multiline = false
	} else if l == 1 {
		if text, ok := t.fragment.children[0].(*Text); ok {
			multiline = strings.Contains(text.text, "\n")
		}
	}

	childpadding := padding + Padding
	if multiline {
		sw(childpadding)
	}
	t.fragment.write(w, childpadding)
	if multiline {
		sw(padding)
	}
	sw("</")
	sw(t.tag)
	sw(">")
}

// VoidTag is a self closing tag that cannot contain children
type VoidTag struct {
	tag string
	attributes
	*parent
}

// NewVoidTag creates a VoidTag. It requires a tag name and optionally takes
// attributes.
func NewVoidTag(tag string, attrs ...string) *VoidTag {
	t := &VoidTag{
		tag:        tag,
		attributes: newAttributes(attrs),
	}
	t.parent = newParent(t)
	return t
}

// Name returns the tag name
func (t *VoidTag) Name() string { return t.tag }

// Write a VoidTag to a writer
func (t *VoidTag) Write(w io.Writer) {
	t.write(ToStringWriter(w), NewLine)
}

func (t *VoidTag) write(w StringWriter, padding string) {
	sw := w.WriteString
	sw("<")
	sw(t.tag)
	t.attributes.write(w)
	sw(" />")
}
