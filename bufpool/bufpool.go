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

// PutAndCopy returns the buffer to the pool and returns a copy of it's byte
// slice
func PutAndCopy(buf *bytes.Buffer) []byte {
	bs := buf.Bytes()
	cp := make([]byte, len(bs))
	copy(cp, bs)
	buf.Reset()
	pool.Put(buf)
	return cp
}

// PutStr returns a buffer from the pool and returns it's value as a string
func PutStr(buf *bytes.Buffer) string {
	return string(PutAndCopy(buf))
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
