package gothicio

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"text/template"
)

func TestStringWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	var sw StringWriter = buf
	sw.WriteString("testing")
	assert.Equal(t, "testing", buf.String())
}

func TestSumWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	sw := NewSumWriter(buf)
	sw.WriteString("test1")
	sw.WriteRune('-')
	sw.Write([]byte("test2"))
	assert.NoError(t, sw.Err)
	sw.Err = fmt.Errorf("test error")
	sw.WriteString("test3")
	assert.Error(t, sw.Err)
	assert.Equal(t, "test1-test2", buf.String())
	assert.Equal(t, int64(11), sw.Sum)
}

func TestMany(t *testing.T) {
	buf := &bytes.Buffer{}
	wm := WriterToMerge(StringWriterTo("test1"), StringWriterTo("test2"))
	wm = WriterToMerge(wm, StringWriterTo("test3"))

	tos := wm.(WriterTos)
	assert.Len(t, tos, 3)
	n, err := MultiWrite(buf, tos, ":")
	assert.NoError(t, err)
	assert.Equal(t, int64(17), n)
	assert.Equal(t, "test1:test2:test3", buf.String())
}

func TestTemplate(t *testing.T) {
	data := struct {
		Test string
	}{
		Test: "testing",
	}
	tmpl := template.Must(template.New("test").Parse(`base template{{define "core"}}My name is {{.Test}}{{end}}`))
	twt := NewTemplateTo(tmpl, "", data)
	buf := &bytes.Buffer{}
	twt.WriteTo(buf)
	assert.Equal(t, "base template", buf.String())
	buf.Reset()
	twt.Name = "core"
	twt.WriteTo(buf)
	assert.Equal(t, "My name is testing", buf.String())
}

func TestWriteTo(t *testing.T) {
	wt := StringWriterTo("this is a test")
	b, err := WriteTo(wt)
	assert.NoError(t, err)
	assert.Equal(t, string(wt), string(b))
}
