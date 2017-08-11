package gothichtml

import (
	"io"
)

// Doctype represents a doctype tag. It does not fulfil TagNode.
type Doctype struct {
	doctype string
	*parent
}

// NewDoctype returns a new doctype tag. To create a standard doctype tag, just
// call NewDoctype("html")
func NewDoctype(doctype string) *Doctype {
	d := &Doctype{
		doctype: doctype,
	}
	d.parent = newParent(d)
	return d
}

// Write Doctype to an io.Writer
func (d *Doctype) Write(w io.Writer) {
	d.write(ToStringWriter(w), NewLine)
}

func (d *Doctype) write(sw StringWriter, Padding string) {
	sw.WriteString("<!DOCTYPE ")
	sw.WriteString(d.doctype)
	sw.WriteString(">")
}
