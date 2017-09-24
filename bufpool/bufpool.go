package bufpool

import (
	"bytes"
	"io"
	"sync"
)

var pool = &sync.Pool{
	New: func() interface{} { return &bytes.Buffer{} },
}

// Get returns a Buffer from the pool
func Get() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

// Put returns a buffer from the pool
func Put(buf *bytes.Buffer) {
	buf.Reset()
	pool.Put(buf)
}

// Put returns a buffer from the pool and returns it's value as a string
func PutStr(buf *bytes.Buffer) string {
	str := buf.String()
	buf.Reset()
	pool.Put(buf)
	return str
}

// TemplateExecutor is an interface representing the ExecuteTemplate method on
// a template.
type TemplateExecutor interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

// ExecuteTemplate using a buffer from the pool.
func ExecuteTemplate(templates TemplateExecutor, name string, data interface{}) (string, error) {
	buf := Get()
	err := templates.ExecuteTemplate(buf, name, data)
	str := buf.String()
	Put(buf)
	return str, err
}
