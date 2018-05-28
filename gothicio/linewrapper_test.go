package gothicio

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicWrapping(t *testing.T) {
	text := []byte("aggrandize epistolography playwoman unreformable wretched supinate reassort relent kurchicine lithophyllous trilingual inventiveness historicoprophetic Bereshith musal unempty Lagothrix symbological zechin soundlessly arylate fetterbush probationism pluriseptate")

	buf := &bytes.Buffer{}
	w := NewLineWrappingWriter(buf)
	w.Write(text)
	expected := `aggrandize epistolography playwoman unreformable wretched supinate reassort
relent kurchicine lithophyllous trilingual inventiveness historicoprophetic Bereshith
musal unempty Lagothrix symbological zechin soundlessly arylate fetterbush
probationism pluriseptate`
	assert.Equal(t, expected, buf.String())

	buf.Reset()
	WrapWidth = 25
	w = NewLineWrappingWriter(buf)
	w.Write(text)
	WrapWidth = 80
	expected = `aggrandize epistolography
playwoman unreformable
wretched supinate reassort
relent kurchicine
lithophyllous trilingual
inventiveness historicoprophetic
Bereshith musal unempty
Lagothrix symbological zechin
soundlessly arylate
fetterbush probationism
pluriseptate`
	assert.Equal(t, expected, buf.String())

	buf.Reset()
	w = NewLineWrappingWriter(
		LineWrapperContextWriter{
			Writer: buf,
			Width:  35,
			Pad:    "// ",
		},
	)
	w.WritePadding()
	w.Write(text)
	expected = `// aggrandize epistolography
// playwoman unreformable wretched
// supinate reassort relent kurchicine
// lithophyllous trilingual
// inventiveness historicoprophetic Bereshith
// musal unempty Lagothrix
// symbological zechin soundlessly arylate
// fetterbush probationism pluriseptate`
	assert.Equal(t, expected, buf.String())
}

func TestContinueWrapping(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewLineWrappingWriter(
		LineWrapperContextWriter{
			Writer: buf,
			Width:  80,
		},
	)
	w.Write([]byte("aggrandize epistolography playwoman unreformable"))
	w.Write([]byte(" wretched supinate reassort relent kurchicine lithophyllous"))
	expected := `aggrandize epistolography playwoman unreformable wretched supinate reassort
relent kurchicine lithophyllous`
	assert.Equal(t, expected, buf.String())
}
