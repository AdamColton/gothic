package gothichtml

import (
	"io"
)

// Text represents an html text node.
type Text struct {
	text string
	*parent
}

// NewText takes the text string and returns a Text Node.
func NewText(text string) *Text {
	t := &Text{
		text: text,
	}
	t.parent = newParent(t)
	return t
}

// Write a Text Node to a writer.
func (t *Text) Write(w io.Writer) {
	t.write(ToStringWriter(w), NewLine)
}
func (t *Text) write(sw StringWriter, Padding string) {
	sw.WriteString(t.text)
}
