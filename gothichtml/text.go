package gothichtml

import (
	"io"
	"strings"
)

var Wrapwidth = 80

// Text represents an html text node.
type Text struct {
	Text string
	Wrap bool
	*parent
}

// NewText takes the text string and returns a Text Node.
func NewText(text string) *Text {
	t := &Text{
		Text: text,
		Wrap: true,
	}
	t.parent = newParent(t)
	return t
}

// Write a Text Node to a writer.
func (t *Text) Write(w io.Writer) {
	t.write(newWriter(w))
}

func (t *Text) write(w *writer) {
	if !t.Wrap || len([]rune(t.Text))+len([]rune(w.padding)) < Wrapwidth || NewLine == "" {
		w.write(t.Text)
		return
	}

	l := len([]rune(w.padding))
	sum := l
	if !w.onNewLine {
		w.nl()
	}
	strs := strings.Split(t.Text, " ")
	last := len(strs) - 1
	for i, str := range strs {
		if str == "" {
			continue
		}
		ln := len([]rune(str))
		sum += ln
		if sum < Wrapwidth {
			if !w.onNewLine {
				sum++
				w.write(" ")
			}
			w.write(str)
		} else if i != last {
			w.nl()
			w.write(str)
			sum = l + ln
		}
	}
	w.pnl()
}
