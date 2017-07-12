package lensnippet

import (
	"fmt"
	"github.com/adamcolton/gothic"
)

const (
	minMode byte = iota
	maxMode
	equalMode
)

type lengthSnippet struct {
	length     int
	errSnippet gothic.Snippet
	on         string
	mode       byte
}

func Min(length int, errSnippet gothic.Snippet) gothic.Snippet {
	return lengthSnippet{
		length:     length,
		errSnippet: errSnippet,
		mode:       minMode,
	}
}

func (ls lengthSnippet) AddContext(key, value string) gothic.Snippet {
	if key == "On" {
		ls.on = value
	}
	return ls
}

const (
	errStr = `%s must be at least %d long`
	minStr = `if len(%s) < %d {
	%s
}`
)

func (ls lengthSnippet) String() string {
	es := ls.errSnippet.AddContext("Error", fmt.Sprintf(errStr, ls.on, ls.length))
	return fmt.Sprintf(minStr, ls.on, ls.length, es)
}
