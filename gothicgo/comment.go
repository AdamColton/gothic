package gothicgo

import (
	"github.com/adamcolton/gothic/gothicio"
	"io"
)

// Comment string that automatically wraps the string
type Comment struct {
	Comment string
	Width   int
}

// NewComment with default CommentWidth
func NewComment(comment string) Comment {
	return Comment{
		Comment: comment,
		Width:   CommentWidth,
	}
}

var nl = []byte("\n")

// WriteTo wraps the comment and writes it to the Writer
func (c Comment) WriteTo(w io.Writer) (int64, error) {
	lww := gothicio.NewLineWrappingWriter(
		gothicio.LineWrapperContextWriter{
			Writer: w,
			Width:  c.Width,
			Pad:    "// ",
		},
	)
	s0 := lww.Sum
	lww.WritePadding()
	lww.Write([]byte(c.Comment))
	lww.SumWriter.Write(nl)

	return lww.Sum - s0, lww.Err
}
